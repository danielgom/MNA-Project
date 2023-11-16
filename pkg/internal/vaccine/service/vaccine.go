package service

import (
	"MNA-project/pkg/internal/pet"
	petservice "MNA-project/pkg/internal/pet/service"
	"MNA-project/pkg/internal/vaccine"
	"MNA-project/pkg/internal/vaccine/db"
	"MNA-project/pkg/internal/vaccine/model"
	"MNA-project/pkg/util/errors"
	"context"
	"github.com/samber/lo"
	"time"
)

type vaccineSvc struct {
	vaccineDb  db.Repository
	petService petservice.PetService
}

func NewVaccineService(uR db.Repository, petService petservice.PetService) VaccineService {
	return &vaccineSvc{vaccineDb: uR, petService: petService}
}

func (v *vaccineSvc) GetById(ctx context.Context, vaccineID, petID, userID int64) (*vaccine.GeneralResponse, errors.CommonError) {
	res, commonError := v.vaccineDb.FindByID(ctx, vaccineID, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return buildVaccineResponse(res, p), nil
}

func (v *vaccineSvc) GetAllByPet(ctx context.Context, petID, userID int64) ([]*vaccine.GeneralResponse, errors.CommonError) {
	vacs, commonError := v.vaccineDb.FindAllByPet(ctx, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return lo.Map(vacs, func(v *model.Vaccine, _ int) *vaccine.GeneralResponse {
		return buildVaccineResponse(v, p)
	}), nil
}

func (v *vaccineSvc) Save(ctx context.Context, req *vaccine.RegisterRequest, petID, userID int64) (*vaccine.GeneralResponse, errors.CommonError) {
	localTime := time.Now().Local()
	t := time.Time(*req.Date)
	nDate := time.Time(*req.NextDate)

	newVaccine := &model.Vaccine{
		PetID:     petID,
		Type:      req.Type,
		VetName:   req.VetName,
		Address:   req.Address,
		Date:      &t,
		NextDate:  &nDate,
		UpdatedAt: &localTime,
	}

	res, commonError := v.vaccineDb.Save(ctx, newVaccine)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return buildVaccineResponse(res, p), commonError
}

func (v *vaccineSvc) UpdateById(ctx context.Context, vaccineID int64, request *vaccine.UpdateRequest) (*vaccine.GeneralResponse, errors.CommonError) {
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

	dbErr := v.vaccineDb.Update(ctx, &model.Vaccine{
		ID:        vaccineID,
		Type:      request.Type,
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

func (v *vaccineSvc) DeleteById(ctx context.Context, vaccineID, petID int64) errors.CommonError {
	return v.vaccineDb.Delete(ctx, vaccineID, petID)
}

func buildVaccineResponse(res *model.Vaccine, pet *pet.GeneralResponse) *vaccine.GeneralResponse {
	d := vaccine.CommonDate(*res.Date)
	nD := vaccine.CommonDate(*res.NextDate)

	return &vaccine.GeneralResponse{
		ID:       res.ID,
		PetName:  pet.Name,
		PetID:    pet.ID,
		VetName:  res.VetName,
		Type:     res.Type,
		Address:  res.Address,
		Date:     &d,
		NextDate: &nD,
	}
}
