package server

import (
	"time"
)

type UserDetail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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
	SessionID             string     `json:"session_id"`
	AccessToken           string     `json:"access_token"`
	AccessTokenExpiresAt  time.Time  `json:"access_token_expires_at"`
	RefreshToken          string     `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time  `json:"refresh_token_expires_at"`
	User                  UserDetail `json:"user"`
}

type CreateTaskRequest struct {
	Title       string `form:"title" json:"title" binding:"required"`
	Description string `form:"description" json:"description"`
	DueDateTime string `form:"due_date_time" json:"due_date_time"`
}

type CreateTaskResponse struct {
	ID string `json:"id"`
}

type ListTasksRequest struct {
	Limit   int  `form:"limit"`
	Offset  int  `form:"offset"`
	Pending bool `form:"pending"`
}

type TaskDetails struct {
	Title       string
	Description string
	Status      string
	DueDate     string
	CreatedAt   string
	UpdatedAt   string
}

type GetTaskResponse struct {
	Task TaskDetails
}

type ListTasksResponse struct {
	Count int            `json:"count"`
	Tasks []*TaskDetails `json:"tasks"`
}

type UpdateTaskRequest struct {
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
	DueDateTime string `form:"due_date_time" json:"due_date_time"`
}

type RenewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}
