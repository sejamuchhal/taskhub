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
		if err == gorm.ErrRecordNotFound {
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

	expiry := time.Now().Add(24 * time.Hour)
	token, err := s.TokenHandler.CreateToken(user.ID, user.Email, expiry)
	if err != nil {
		logger.WithError(err).Error("Error creating access token")
		return nil, status.Errorf(codes.Internal, "Error creating access token: %v", err)
	}

	res := &pb.LoginResponse{
		Token: token,
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

    claims, err := s.TokenHandler.VerifyToken(req.Token)
    if err != nil {
        logger.WithError(err).Error("Invalid token")
        return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
    }

    return &pb.ValidateResponse{
        UserId: claims.UserID,
        Email:  claims.Email,
    }, nil
}