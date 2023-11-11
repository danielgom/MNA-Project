// Package db contains all repositories used by this API.
package db

import (
	"context"

	"MNA-project/pkg/internal/user/model"
	"MNA-project/pkg/util/errors"
)

// Repository serves as a middleware to call our users table.
type Repository interface {
	FindByEmail(context.Context, string) (*model.User, errors.CommonError)
	FindById(context.Context, int64) (*model.User, errors.CommonError)
	FindAll(context.Context) ([]*model.User, errors.CommonError)
	Save(context.Context, *model.User) (*model.User, errors.CommonError)
	Update(context.Context, *model.User) errors.CommonError
	Delete(context.Context, int64) errors.CommonError
}
