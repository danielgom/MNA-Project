// Package service contains all core logic for this API.
package service

import (
	"context"

	"MNA-project/pkg/internal/user"
	"MNA-project/pkg/util/errors"
)

// UserService contains all the business logic for the user.
type UserService interface {
	SignUp(context.Context, *user.RegisterRequest) (*user.RegisterResponse, errors.CommonError)
	GetById(context.Context, int64) (*user.Response, errors.CommonError)
	GetAll(context.Context) ([]*user.Response, errors.CommonError)
	UpdateById(context.Context, int64, *user.UpdateRequest) (*user.Response, errors.CommonError)
	DeleteById(context.Context, int64) errors.CommonError
	Login(context.Context, *user.LoginRequest) (*user.LoginResponse, errors.CommonError)
}
