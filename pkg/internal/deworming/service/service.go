package service

import (
	"MNA-project/pkg/internal/deworming"
	"context"

	"MNA-project/pkg/util/errors"
)

type DewormingService interface {
	GetById(context.Context, int64, int64, int64) (*deworming.GeneralResponse, errors.CommonError)
	GetAllByPet(context.Context, int64, int64) ([]*deworming.GeneralResponse, errors.CommonError)
	UpdateById(context.Context, int64, *deworming.UpdateRequest) (*deworming.GeneralResponse, errors.CommonError)
	Save(context.Context, *deworming.RegisterRequest, int64, int64) (*deworming.GeneralResponse, errors.CommonError)
	DeleteById(context.Context, int64, int64) errors.CommonError
}
