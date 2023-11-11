// Package user contains structs and components needed for the user domain
package user

import (
	"MNA-project/pkg/internal/user/model"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// RegisterRequest comes from the signup request.
type RegisterRequest struct {
	Password string `json:"password" validate:"required,password"`
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	LastName string `json:"last_name" validate:"required"`
}

type UpdateRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
}

type UpdatePasswordRequest struct {
	Password string `json:"password"`
}

// RegisterResponse is the struct for a successful signUp.
type RegisterResponse struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	LastName  string     `json:"last_name"`
	FullName  string     `json:"full_name"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// Response is a struct with the user information.
type Response struct {
	ID        int64      `json:"id"`
	Email     string     `json:"email"`
	Name      string     `json:"name"`
	LastName  string     `json:"last_name"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Name      string     `json:"name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Token     string     `json:"token"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type Users struct {
	Users []*Response `json:"users"`
}

// BuildRegisterResponse builds the output of the signUp response when is not error -ed.
func BuildRegisterResponse(user *model.User) *RegisterResponse {
	var response RegisterResponse

	response.ID = user.ID
	response.Email = user.Email
	response.CreatedAt = user.CreatedAt
	response.Name = user.Name
	response.LastName = user.LastName
	response.UpdatedAt = user.UpdatedAt
	response.FullName = user.Name + " " + user.LastName

	return &response
}

func BuildLoginResponse(user *model.User, expTime *time.Time, token string) *LoginResponse {
	return &LoginResponse{
		Name:      user.Name,
		LastName:  user.LastName,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: expTime,
	}
}
