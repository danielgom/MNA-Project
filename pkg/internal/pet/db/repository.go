package db

import (
	"context"

	"MNA-project/pkg/internal/pet/model"
	"MNA-project/pkg/util/errors"
)

type Repository interface {
	FindByID(context.Context, int64, int64) (*model.Pet, errors.CommonError)
	FindAllByUser(context.Context, int64) ([]*model.Pet, errors.CommonError)
	Save(context.Context, *model.Pet) (*model.Pet, errors.CommonError)
	Update(context.Context, *model.Pet) errors.CommonError
	Delete(context.Context, int64, int64) errors.CommonError
}
