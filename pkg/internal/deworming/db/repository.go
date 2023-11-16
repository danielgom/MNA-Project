package db

import (
	"MNA-project/pkg/internal/deworming/model"
	"MNA-project/pkg/util/errors"
	"context"
)

type Repository interface {
	FindByID(context.Context, int64, int64) (*model.Deworming, errors.CommonError)
	FindAllByPet(context.Context, int64) ([]*model.Deworming, errors.CommonError)
	Save(context.Context, *model.Deworming) (*model.Deworming, errors.CommonError)
	Update(context.Context, *model.Deworming) errors.CommonError
	Delete(context.Context, int64, int64) errors.CommonError
}
