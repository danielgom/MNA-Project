package service

import (
	"MNA-project/pkg/internal/pet"
	petservice "MNA-project/pkg/internal/pet/service"
	"MNA-project/pkg/internal/surgeries"
	"MNA-project/pkg/internal/surgeries/db"
	"MNA-project/pkg/internal/surgeries/model"
	"MNA-project/pkg/util/errors"
	"context"
	"github.com/samber/lo"
	"time"
)

type surgerySvc struct {
	surgeryDb  db.Repository
	petService petservice.PetService
}

func NewSurgeryService(uR db.Repository, petService petservice.PetService) SurgeryService {
	return &surgerySvc{surgeryDb: uR, petService: petService}
}

func (v *surgerySvc) GetById(ctx context.Context, surgeryID, petID, userID int64) (*surgeries.GeneralResponse, errors.CommonError) {
	res, commonError := v.surgeryDb.FindByID(ctx, surgeryID, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("vet visit not found")
	}

	return buildSurgeryResponse(res, p), nil
}

func (v *surgerySvc) GetAllByPet(ctx context.Context, petID, userID int64) ([]*surgeries.GeneralResponse, errors.CommonError) {
	vacs, commonError := v.surgeryDb.FindAllByPet(ctx, petID)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return lo.Map(vacs, func(v *model.Surgery, _ int) *surgeries.GeneralResponse {
		return buildSurgeryResponse(v, p)
	}), nil
}

func (v *surgerySvc) Save(ctx context.Context, req *surgeries.RegisterRequest, petID, userID int64) (*surgeries.GeneralResponse, errors.CommonError) {
	localTime := time.Now().Local()
	t := time.Time(*req.Date)

	newVaccine := &model.Surgery{
		PetID:     petID,
		VetName:   req.VetName,
		Address:   req.Address,
		Name:      req.Name,
		Comments:  req.Comments,
		Date:      &t,
		UpdatedAt: &localTime,
	}

	res, commonError := v.surgeryDb.Save(ctx, newVaccine)
	if commonError != nil {
		return nil, commonError
	}

	p, userErr := v.petService.GetById(ctx, petID, userID)
	if userErr != nil {
		return nil, errors.NewNotFoundError("pet not found")
	}

	return buildSurgeryResponse(res, p), commonError
}

func (v *surgerySvc) UpdateById(ctx context.Context, surgeryID int64, request *surgeries.UpdateRequest) (*surgeries.GeneralResponse, errors.CommonError) {
	updateTime := time.Now().Local()

	var t *time.Time
	if request.Date == nil {
		t = nil
	} else {
		tmp := time.Time(*request.Date)
		t = &tmp
	}

	dbErr := v.surgeryDb.Update(ctx, &model.Surgery{
		ID:        surgeryID,
		VetName:   request.VetName,
		Address:   request.Address,
		Name:      request.Name,
		Comments:  request.Comments,
		Date:      t,
		UpdatedAt: &updateTime,
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return v.GetById(ctx, surgeryID, request.PetID, request.UserID)
}

func (v *surgerySvc) DeleteById(ctx context.Context, surgeryID, petID int64) errors.CommonError {
	return v.surgeryDb.Delete(ctx, surgeryID, petID)
}

func buildSurgeryResponse(res *model.Surgery, pet *pet.GeneralResponse) *surgeries.GeneralResponse {
	d := surgeries.CommonDate(*res.Date)

	return &surgeries.GeneralResponse{
		ID:       res.ID,
		PetName:  pet.Name,
		PetID:    pet.ID,
		VetName:  res.VetName,
		Address:  res.Address,
		Date:     &d,
		Name:     res.Name,
		Comments: res.Comments,
	}
}
