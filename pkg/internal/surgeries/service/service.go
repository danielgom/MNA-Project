package service

import (
	"MNA-project/pkg/internal/surgeries"
	"context"

	"MNA-project/pkg/util/errors"
)

type SurgeryService interface {
	GetById(context.Context, int64, int64, int64) (*surgeries.GeneralResponse, errors.CommonError)
	GetAllByPet(context.Context, int64, int64) ([]*surgeries.GeneralResponse, errors.CommonError)
	UpdateById(context.Context, int64, *surgeries.UpdateRequest) (*surgeries.GeneralResponse, errors.CommonError)
	Save(context.Context, *surgeries.RegisterRequest, int64, int64) (*surgeries.GeneralResponse, errors.CommonError)
	DeleteById(context.Context, int64, int64) errors.CommonError
}
