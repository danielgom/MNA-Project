package db

import (
	"MNA-project/pkg/internal/surgeries/model"
	"MNA-project/pkg/util/errors"
	"context"
)

type Repository interface {
	FindByID(context.Context, int64, int64) (*model.Surgery, errors.CommonError)
	FindAllByPet(context.Context, int64) ([]*model.Surgery, errors.CommonError)
	Save(context.Context, *model.Surgery) (*model.Surgery, errors.CommonError)
	Update(context.Context, *model.Surgery) errors.CommonError
	Delete(context.Context, int64, int64) errors.CommonError
}
