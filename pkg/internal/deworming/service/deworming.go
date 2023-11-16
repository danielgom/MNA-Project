package service

import (
	"MNA-project/pkg/internal/deworming"
	"MNA-project/pkg/internal/deworming/db"
	"MNA-project/pkg/internal/deworming/model"
	"MNA-project/pkg/internal/pet"
	petservice "MNA-project/pkg/internal/pet/service"
	"MNA-project/pkg/util/errors"
	"context"
	"github.com/samber/lo"
	"time"
)

type dewormingSvc struct {
	dewormingDb db.Repository
	petService  petservice.PetService
}

func NewDewormingService(uR db.Repository, petService petservice.PetService) DewormingService {
	return &dewormingSvc{dewormingDb: uR, petService: petService}
}

func (v *dewormingSvc) GetById(ctx context.Context, vaccineID, petID, userID int64) (*deworming.GeneralResponse, errors.CommonError) {
	res, commonError := v.dewormingDb.FindByID(ctx, vaccineID, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return buildVaccineResponse(res, p), nil
}

func (v *dewormingSvc) GetAllByPet(ctx context.Context, petID, userID int64) ([]*deworming.GeneralResponse, errors.CommonError) {
	vacs, commonError := v.dewormingDb.FindAllByPet(ctx, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return lo.Map(vacs, func(v *model.Deworming, _ int) *deworming.GeneralResponse {
		return buildVaccineResponse(v, p)
	}), nil
}

func (v *dewormingSvc) Save(ctx context.Context, req *deworming.RegisterRequest, petID, userID int64) (*deworming.GeneralResponse, errors.CommonError) {
	localTime := time.Now().Local()
	t := time.Time(*req.Date)
	nDate := time.Time(*req.NextDate)

	newVaccine := &model.Deworming{
		PetID:     petID,
		VetName:   req.VetName,
		Address:   req.Address,
		Date:      &t,
		NextDate:  &nDate,
		UpdatedAt: &localTime,
	}

	res, commonError := v.dewormingDb.Save(ctx, newVaccine)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return buildVaccineResponse(res, p), commonError
}

func (v *dewormingSvc) UpdateById(ctx context.Context, vaccineID int64, request *deworming.UpdateRequest) (*deworming.GeneralResponse, errors.CommonError) {
	updateTime := time.Now().Local()

	var t *time.Time
	if request.Date == nil {
		t = nil
	} else {
		tmp := time.Time(*request.Date)
		t = &tmp
	}

	var nDate *time.Time
	if request.NextDate == nil {
		nDate = nil
	} else {
		tmp := time.Time(*request.NextDate)
		nDate = &tmp
	}

	dbErr := v.dewormingDb.Update(ctx, &model.Deworming{
		ID:        vaccineID,
		VetName:   request.VetName,
		Address:   request.Address,
		Date:      t,
		NextDate:  nDate,
		UpdatedAt: &updateTime,
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return v.GetById(ctx, vaccineID, request.PetID, request.UserID)
}

func (v *dewormingSvc) DeleteById(ctx context.Context, vaccineID, petID int64) errors.CommonError {
	return v.dewormingDb.Delete(ctx, vaccineID, petID)
}

func buildVaccineResponse(res *model.Deworming, pet *pet.GeneralResponse) *deworming.GeneralResponse {
	d := deworming.CommonDate(*res.Date)
	nD := deworming.CommonDate(*res.NextDate)

	return &deworming.GeneralResponse{
		ID:       res.ID,
		PetName:  pet.Name,
		PetID:    pet.ID,
		VetName:  res.VetName,
		Address:  res.Address,
		Date:     &d,
		NextDate: &nD,
	}
}
