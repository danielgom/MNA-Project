package db

import (
	"MNA-project/pkg/internal/vet_visits/model"
	"MNA-project/pkg/util/errors"
	"context"
)

type Repository interface {
	FindByID(context.Context, int64, int64) (*model.VetVisit, errors.CommonError)
	FindAllByPet(context.Context, int64) ([]*model.VetVisit, errors.CommonError)
	Save(context.Context, *model.VetVisit) (*model.VetVisit, errors.CommonError)
	Update(context.Context, *model.VetVisit) errors.CommonError
	Delete(context.Context, int64, int64) errors.CommonError
}
