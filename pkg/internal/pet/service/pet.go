package service

import (
	"MNA-project/pkg/internal/user"
	"context"
	"github.com/samber/lo"
	"time"

	"MNA-project/pkg/internal/pet"
	"MNA-project/pkg/internal/pet/db"
	pm "MNA-project/pkg/internal/pet/model"
	"MNA-project/pkg/internal/user/service"
	"MNA-project/pkg/util/errors"
)

type petSvc struct {
	petDB       db.Repository
	userService service.UserService
}

// NewPetService returns a new instance of pet service.
func NewPetService(uR db.Repository, us service.UserService) PetService {
	return &petSvc{petDB: uR, userService: us}
}

func (p *petSvc) GetById(ctx context.Context, petID int64, userID int64) (*pet.GeneralResponse, errors.CommonError) {
	res, commonError := p.petDB.FindByID(ctx, petID, userID)
	if commonError != nil {
		return nil, commonError
	}

	u, userErr := p.userService.GetById(ctx, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	return buildPetResponse(res, u), nil
}

func (p *petSvc) GetAllByUser(ctx context.Context, userID int64) ([]*pet.GeneralResponse, errors.CommonError) {
	pets, commonError := p.petDB.FindAllByUser(ctx, userID)
	if commonError != nil {
		return nil, commonError
	}

	u, userErr := p.userService.GetById(ctx, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	return lo.Map(pets, func(p *pm.Pet, _ int) *pet.GeneralResponse {
		return buildPetResponse(p, u)
	}), nil
}

func (p *petSvc) Save(ctx context.Context, req *pet.RegisterRequest, userID int64) (*pet.GeneralResponse, errors.CommonError) {
	localTime := time.Now().Local()
	t := time.Time(req.BirthDate)

	newPet := &pm.Pet{
		UserID:       userID,
		Name:         req.Name,
		Age:          req.Age,
		Breed:        req.Breed,
		BirthDate:    &t,
		RegisterDate: &localTime,
		UpdatedAt:    &localTime,
	}

	res, commonError := p.petDB.Save(ctx, newPet)
	if commonError != nil {
		return nil, commonError
	}

	u, userErr := p.userService.GetById(ctx, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("user not found")
	}

	return buildPetResponse(res, u), commonError
}

func (p *petSvc) UpdateById(ctx context.Context, petID int64, request *pet.UpdateRequest) (*pet.GeneralResponse, errors.CommonError) {
	updateTime := time.Now().Local()
	dbErr := p.petDB.Update(ctx, &pm.Pet{
		ID:        petID,
		Name:      request.Name,
		Age:       request.Age,
		Breed:     request.Breed,
		UpdatedAt: &updateTime,
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return p.GetById(ctx, petID, request.UserID)
}

func (p *petSvc) DeleteById(ctx context.Context, petID int64, userID int64) errors.CommonError {
	return p.petDB.Delete(ctx, petID, userID)
}

func buildPetResponse(p *pm.Pet, u *user.Response) *pet.GeneralResponse {
	t := pet.BirthDate(*p.BirthDate)
	return &pet.GeneralResponse{
		ID:           p.ID,
		OwnerName:    u.Name,
		OwnerID:      u.ID,
		Name:         p.Name,
		Age:          p.Age,
		Breed:        p.Breed,
		BirthDate:    t,
		RegisterDate: p.RegisterDate,
		UpdatedAt:    p.UpdatedAt,
	}
}
