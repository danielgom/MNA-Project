package service

import (
	"context"

	"MNA-project/pkg/internal/pet"
	"MNA-project/pkg/util/errors"
)

// PetService contains all the business logic for the pet.
type PetService interface {
	GetById(context.Context, int64, int64) (*pet.GeneralResponse, errors.CommonError)
	GetAllByUser(context.Context, int64) ([]*pet.GeneralResponse, errors.CommonError)
	UpdateById(context.Context, int64, *pet.UpdateRequest) (*pet.GeneralResponse, errors.CommonError)
	Save(context.Context, *pet.RegisterRequest, int64) (*pet.GeneralResponse, errors.CommonError)
	DeleteById(context.Context, int64, int64) errors.CommonError
}
