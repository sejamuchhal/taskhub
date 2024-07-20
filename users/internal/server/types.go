package server

import (
	"time"
)

type userDetail struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type SignupUserRequest struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=100"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=6,max=100"`
}

type SignupUserResponse struct {
	ID string `json:"id"`
}

type LoginUserRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken          string     `json:"access_token"`
	AccessTokenExpiresAt time.Time  `json:"access_token_expires_at"`
	User                 userDetail `json:"user"`
}
