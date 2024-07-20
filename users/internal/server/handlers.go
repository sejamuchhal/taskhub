package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/sejamuchhal/task-management/users/internal/database"
	"github.com/sejamuchhal/task-management/users/utils"
)

func (s *Server) Health(c *gin.Context) {
	resp := map[string]string{
		"message": "It's healthy",
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) SignupUser(c *gin.Context) {
	s.logger.Info("Incoming signup request")
	var req SignupUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.WithError(err).Error("Error parsing signup request payload")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.WithError(err).Error("Error hashing password")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = s.db.CreateUser(&database.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				s.logger.WithError(err).Error("User already exists")
				c.JSON(http.StatusForbidden, errorResponse(errors.New("user already exists")))
				return
			}
		}
		s.logger.WithError(err).Error("Error creating user in the database")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	s.logger.Info("User signup successful")
	c.JSON(http.StatusOK, gin.H{"message": "User signup successful"})
}

func (s *Server) LoginUser(c *gin.Context) {
	s.logger.Info("Incoming login request")
	var req LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.WithError(err).Error("Error parsing login request payload")
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.db.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.WithError(err).Warn("User not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid email or password"})
			return
		}
		s.logger.WithError(err).Error("Error fetching user from the database")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		s.logger.WithError(err).Warn("Invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	expiry := time.Now().Add(time.Hour * 24)

	token, err := s.tokenHandler.CreateToken(user.Email, expiry)
	if err != nil {
		s.logger.WithError(err).Error("Error creating access token")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := LoginUserResponse{
		AccessToken:          token,
		AccessTokenExpiresAt: expiry,
		User: userDetail{
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}
	s.logger.Info("User login successful")
	c.JSON(http.StatusOK, res)
}

func (s *Server) CreateTask(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve email from context"})
		return
	}
	user, err := s.db.GetUserByEmail(email.(string))
	if err != nil {
		s.logger.WithError(err).Error("Error fetching user from the database")
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	

	c.JSON(http.StatusOK, gin.H{
		"email": email,
	})
}
