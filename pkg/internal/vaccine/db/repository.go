package db

import (
	"MNA-project/pkg/internal/vaccine/model"
	"MNA-project/pkg/util/errors"
	"context"
)

type Repository interface {
	FindByID(context.Context, int64, int64) (*model.Vaccine, errors.CommonError)
	FindAllByPet(context.Context, int64) ([]*model.Vaccine, errors.CommonError)
	Save(context.Context, *model.Vaccine) (*model.Vaccine, errors.CommonError)
	Update(context.Context, *model.Vaccine) errors.CommonError
	Delete(context.Context, int64, int64) errors.CommonError
}
