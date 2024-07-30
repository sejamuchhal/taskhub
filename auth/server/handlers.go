package server

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	pb "github.com/sejamuchhal/taskhub/auth/pb"
	"github.com/sejamuchhal/taskhub/auth/storage"
	"github.com/sejamuchhal/taskhub/auth/util"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

// Signup handles user registration
// Checks if user already exist for email or not if yes, retrun error
func (s *Server) Signup(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"method": "Signup",
		"req":    req,
	})
	logger.Debug("Incoming signup request")

	user, err := s.Storage.GetUserByEmail(req.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.WithError(err).Error("Error fetching user from the database")
			return nil, status.Errorf(codes.Internal, "Error fetching user from the database: %v", err)
		}
	} else if user != nil {
		logger.WithError(err).Error("User already exists")
		return nil, status.Errorf(codes.AlreadyExists, "User already exists")
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		logger.WithError(err).Error("Error hashing password")
		return nil, status.Errorf(codes.Internal, "Error hashing password: %v", err)
	}

	err = s.Storage.CreateUser(&storage.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			logger.WithError(err).Error("User already exists")
			return nil, status.Errorf(codes.AlreadyExists, "User already exists")
		}
		logger.WithError(err).Error("Error creating user in the database")
		return nil, status.Errorf(codes.Internal, "Error creating user in the database: %v", err)
	}

	logger.Debug("User signup successful")
	return &pb.SignupResponse{Message: "User signup successful"}, nil
}

// Login handles user login
// Check if user exist or not for email
// if yes, compare password hash
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"method": "Login",
		"req":    req,
	})
	logger.Debug("Incoming login request")

	user, err := s.Storage.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithError(err).Warn("User not found")
			return nil, status.Errorf(codes.NotFound, "Invalid email or password")
		}
		logger.WithError(err).Error("Error fetching user from the database")
		return nil, status.Errorf(codes.Internal, "Error fetching user from the database: %v", err)
	}

	if err := util.CheckPasswordHash(req.Password, user.Password); err != nil {
		logger.WithError(err).Warn("Invalid password")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid email or password")
	}

	accessToken, accessClaims, err := s.TokenHandler.CreateToken(user, s.Config.AccessTokenDuration, "access")
	if err != nil {
		logger.WithError(err).Error("Error creating access token")
		return nil, status.Errorf(codes.Internal, "Error creating access token: %v", err)
	}

	refreshToken, refreshClaims, err := s.TokenHandler.CreateToken(user, s.Config.RefreshTokenDuration, "refresh")
	if err != nil {
		logger.WithError(err).Error("Error creating refresh token")
		return nil, status.Errorf(codes.Internal, "Error creating refresh token: %v", err)
	}

	err = s.Storage.CreateSession(&storage.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		Email:        user.Email,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		logger.WithError(err).Error("Error creating session")
		return nil, status.Errorf(codes.Internal, "Error creating session: %v", err)
	}

	res := &pb.LoginResponse{
		SessionId:             refreshClaims.RegisteredClaims.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessClaims.RegisteredClaims.ExpiresAt.Time),
		RefreshTokenExpiresAt: timestamppb.New(refreshClaims.RegisteredClaims.ExpiresAt.Time),
		User: &pb.UserDetail{
			Name:  user.Name,
			Email: user.Email,
		},
	}

	logger.Debug("User login successful")
	return res, nil
}

// Validate token and return
func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	logger := s.Logger.WithField("method", "Validate")
	logger.Debug("Incoming request")

	claims, err := s.TokenHandler.VerifyToken(req.Token, "access")
	if err != nil {
		logger.WithError(err).Error("Invalid token")
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}
	logger.WithField("claims", claims).Debug("Fetched claims")

	return &pb.ValidateResponse{
		UserId: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}

func (s *Server) RenewAccessToken(ctx context.Context, req *pb.RenewAccessTokenRequest) (*pb.RenewAccessTokenResponse, error) {
	logger := s.Logger.WithField("method", "RenewAccessToken")
	logger.Debug("Incoming request")

	refreshClaims, err := s.TokenHandler.VerifyToken(req.RefreshToken, "refresh")
	if err != nil {
		logger.WithError(err).Error("Error veryfying token")
		return nil, status.Errorf(codes.Unauthenticated, "Error veryfying token: %v", err)
	}
	session, err := s.Storage.GetSessionByID(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		logger.WithError(err).Error("Error fethcing session")
		return nil, status.Errorf(codes.Internal, "Error fethcing session: %v", err)

	}

	if session.IsBlocked {
		logger.WithError(err).Error("Session blocked")
		return nil, status.Error(codes.Unauthenticated, "Session blocked")
	}
	if session.Email != refreshClaims.Email || session.RefreshToken != req.RefreshToken {
		logger.WithError(err).Error("Invalid session")
		return nil, status.Error(codes.Unauthenticated, "Invalid session")
	}

	if time.Now().After(session.ExpiresAt) {
		logger.WithError(err).Error("Expired session")
		return nil, status.Error(codes.Unauthenticated, "Expired session")
	}

	accessToken, accessClaims, err := s.TokenHandler.CreateToken(&storage.User{
		ID:    refreshClaims.UserID,
		Email: refreshClaims.Email,
		Role:  refreshClaims.Role,
	}, s.Config.AccessTokenDuration, "access")

	if err != nil {
		logger.WithError(err).Error("Error creating token")
		return nil, status.Errorf(codes.Internal, "Error creating token: %v", err)
	}
	res := &pb.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessClaims.RegisteredClaims.ExpiresAt.Time),
	}
	return res, nil

}

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	logger := s.Logger.WithFields(logrus.Fields{
		"method": "Logout",
		"req":    req,
	})
	logger.Debug("Incoming logout request")

	_, err := s.TokenHandler.VerifyToken(req.AccessToken, "access")
	if err != nil {
		logger.WithError(err).Error("Invalid access tken")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid access token")
	}

	err = s.TokenHandler.BlacklistToken(req.AccessToken, s.Config.AccessTokenDuration)
	if err != nil {
		logger.WithError(err).Error("Error blacklisting access token")
		return nil, status.Errorf(codes.Internal, "Error logging user out")
	}

	refreshClaims, err := s.TokenHandler.VerifyToken(req.RefreshToken, "refresh")
	if err != nil {
		logger.WithError(err).Error("Invalid refresh token")
		return nil, status.Errorf(codes.InvalidArgument, "Invalid refresh token")
	}

	session, err := s.Storage.GetSessionByID(refreshClaims.RegisteredClaims.ID)
	if err != nil {
		logger.WithError(err).Error("Error fetching session")
		return nil, status.Errorf(codes.Internal, "Error fetching session: %v", err)
	}

	if session.IsBlocked {
		logger.Warning("Session already revoked")
		return &pb.LogoutResponse{}, nil
	}

	err = s.Storage.BlockSessionByID(session.ID)
	if err != nil {
		logger.WithError(err).Error("Error blocking session")
		return nil, status.Errorf(codes.Internal, "Error blocking session: %v", err)
	}

	logger.Debug("User logout successful")
	return &pb.LogoutResponse{}, nil
}
