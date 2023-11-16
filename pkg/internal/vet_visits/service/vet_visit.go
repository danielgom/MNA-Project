package service

import (
	"MNA-project/pkg/internal/pet"
	petservice "MNA-project/pkg/internal/pet/service"
	"MNA-project/pkg/internal/vet_visits"
	"MNA-project/pkg/internal/vet_visits/db"
	"MNA-project/pkg/internal/vet_visits/model"
	"MNA-project/pkg/util/errors"
	"context"
	"github.com/samber/lo"
	"time"
)

type vetVisitSvc struct {
	vetVisitDb db.Repository
	petService petservice.PetService
}

func NewVetVisitService(uR db.Repository, petService petservice.PetService) VetVisitService {
	return &vetVisitSvc{vetVisitDb: uR, petService: petService}
}

func (v *vetVisitSvc) GetById(ctx context.Context, vaccineID, petID, userID int64) (*vet_visits.GeneralResponse, errors.CommonError) {
	res, commonError := v.vetVisitDb.FindByID(ctx, vaccineID, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("vet visit not found")
	}

	return buildVetVisitResponse(res, p), nil
}

func (v *vetVisitSvc) GetAllByPet(ctx context.Context, petID, userID int64) ([]*vet_visits.GeneralResponse, errors.CommonError) {
	vacs, commonError := v.vetVisitDb.FindAllByPet(ctx, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return lo.Map(vacs, func(v *model.VetVisit, _ int) *vet_visits.GeneralResponse {
		return buildVetVisitResponse(v, p)
	}), nil
}

func (v *vetVisitSvc) Save(ctx context.Context, req *vet_visits.RegisterRequest, petID, userID int64) (*vet_visits.GeneralResponse, errors.CommonError) {
	localTime := time.Now().Local()
	t := time.Time(*req.Date)

	newVaccine := &model.VetVisit{
		PetID:     petID,
		VetName:   req.VetName,
		Address:   req.Address,
		Reason:    req.Reason,
		Comments:  req.Comments,
		Date:      &t,
		UpdatedAt: &localTime,
	}

	res, commonError := v.vetVisitDb.Save(ctx, newVaccine)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return buildVetVisitResponse(res, p), commonError
}

func (v *vetVisitSvc) UpdateById(ctx context.Context, vaccineID int64, request *vet_visits.UpdateRequest) (*vet_visits.GeneralResponse, errors.CommonError) {
	updateTime := time.Now().Local()

	var t *time.Time
	if request.Date == nil {
		t = nil
	} else {
		tmp := time.Time(*request.Date)
		t = &tmp
	}

	dbErr := v.vetVisitDb.Update(ctx, &model.VetVisit{
		ID:        vaccineID,
		VetName:   request.VetName,
		Address:   request.Address,
		Reason:    request.Reason,
		Comments:  request.Comments,
		Date:      t,
		UpdatedAt: &updateTime,
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return v.GetById(ctx, vaccineID, request.PetID, request.UserID)
}

func (v *vetVisitSvc) DeleteById(ctx context.Context, vaccineID, petID int64) errors.CommonError {
	return v.vetVisitDb.Delete(ctx, vaccineID, petID)
}

func buildVetVisitResponse(res *model.VetVisit, pet *pet.GeneralResponse) *vet_visits.GeneralResponse {
	d := vet_visits.CommonDate(*res.Date)

	return &vet_visits.GeneralResponse{
		ID:       res.ID,
		PetName:  pet.Name,
		PetID:    pet.ID,
		VetName:  res.VetName,
		Address:  res.Address,
		Date:     &d,
		Reason:   res.Reason,
		Comments: res.Comments,
	}
}
