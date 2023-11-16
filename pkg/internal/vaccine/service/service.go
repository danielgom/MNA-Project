package service

import (
	"MNA-project/pkg/internal/vaccine"
	"context"

	"MNA-project/pkg/util/errors"
)

type VaccineService interface {
	GetById(context.Context, int64, int64, int64) (*vaccine.GeneralResponse, errors.CommonError)
	GetAllByPet(context.Context, int64, int64) ([]*vaccine.GeneralResponse, errors.CommonError)
	UpdateById(context.Context, int64, *vaccine.UpdateRequest) (*vaccine.GeneralResponse, errors.CommonError)
	Save(context.Context, *vaccine.RegisterRequest, int64, int64) (*vaccine.GeneralResponse, errors.CommonError)
	DeleteById(context.Context, int64, int64) errors.CommonError
}
