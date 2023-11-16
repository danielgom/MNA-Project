package routes

import (
	"MNA-project/pkg/context"
	"MNA-project/pkg/internal/vaccine"
	"MNA-project/pkg/internal/vaccine/service"
	"MNA-project/pkg/util/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	vacSvc service.VaccineService
}

func NewHandler(svc service.VaccineService) *Handler {
	return &Handler{vacSvc: svc}
}

func (h *Handler) Register(r *echo.Echo, handler func(fn func(context.Context) error) echo.HandlerFunc) {
	authGroup := r.Group("v1/pets")
	authGroup.POST("/:pet_id/vaccines", handler(h.Save))
	authGroup.GET("/:pet_id/vaccines", handler(h.GetVaccinesByPet))
	authGroup.GET("/:pet_id/vaccines/:id", handler(h.GetVaccineByID))
	authGroup.PUT("/:pet_id/vaccines/:id", handler(h.UpdateById))
	authGroup.DELETE("/:pet_id/vaccines/:id", handler(h.DeleteVaccineByID))
}

func (h *Handler) Save(c context.Context) error {
	var req vaccine.RegisterRequest

	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid pet_id"))
	}

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, signErr := h.vacSvc.Save(c.Request().Context(), &req, int64(petID), c.GetUserID())
		if signErr != nil {
			return nil, signErr
		}

		return &context.GResponse{
			Status:   http.StatusCreated,
			Response: res,
		}, nil
	})
}

func (h *Handler) GetVaccinesByPet(c context.Context) error {
	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid pet_id"))
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.vacSvc.GetAllByPet(c.Request().Context(), int64(petID), c.GetUserID())
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: vaccine.Vaccines{Vaccines: res},
		}, nil
	})
}

func (h *Handler) GetVaccineByID(c context.Context) error {
	vaccineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		userID := c.GetUserID()
		res, getErr := h.vacSvc.GetById(c.Request().Context(), int64(vaccineID), int64(petID), userID)
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

func (h *Handler) UpdateById(c context.Context) error {
	var req vaccine.UpdateRequest

	vaccineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		req.UserID = c.GetUserID()
		req.PetID = int64(petID)
		res, delErr := h.vacSvc.UpdateById(c.Request().Context(), int64(vaccineID), &req)
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

func (h *Handler) DeleteVaccineByID(c context.Context) error {
	vaccineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		delErr := h.vacSvc.DeleteById(c.Request().Context(), int64(vaccineID), int64(petID))
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status: http.StatusOK,
		}, nil
	})
}
