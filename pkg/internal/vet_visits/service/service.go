package service

import (
	"MNA-project/pkg/internal/vet_visits"
	"context"

	"MNA-project/pkg/util/errors"
)

type VetVisitService interface {
	GetById(context.Context, int64, int64, int64) (*vet_visits.GeneralResponse, errors.CommonError)
	GetAllByPet(context.Context, int64, int64) ([]*vet_visits.GeneralResponse, errors.CommonError)
	UpdateById(context.Context, int64, *vet_visits.UpdateRequest) (*vet_visits.GeneralResponse, errors.CommonError)
	Save(context.Context, *vet_visits.RegisterRequest, int64, int64) (*vet_visits.GeneralResponse, errors.CommonError)
	DeleteById(context.Context, int64, int64) errors.CommonError
}
