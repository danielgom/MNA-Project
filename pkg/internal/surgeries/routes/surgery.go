package routes

import (
	"MNA-project/pkg/context"
	"MNA-project/pkg/internal/surgeries"
	"MNA-project/pkg/internal/surgeries/service"
	"MNA-project/pkg/util/errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type Handler struct {
	surgerySvc service.SurgeryService
}

func NewHandler(svc service.SurgeryService) *Handler {
	return &Handler{surgerySvc: svc}
}

func (h *Handler) Register(r *echo.Echo, handler func(fn func(context.Context) error) echo.HandlerFunc) {
	authGroup := r.Group("v1/pets")
	authGroup.POST("/:pet_id/surgeries", handler(h.Save))
	authGroup.GET("/:pet_id/surgeries", handler(h.GetSurgeriesByPet))
	authGroup.GET("/:pet_id/surgeries/:id", handler(h.GetSurgeryByID))
	authGroup.PUT("/:pet_id/surgeries/:id", handler(h.UpdateById))
	authGroup.DELETE("/:pet_id/surgeries/:id", handler(h.DeleteVetVisitByID))
}

func (h *Handler) Save(c context.Context) error {
	var req surgeries.RegisterRequest

	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid pet_id"))
	}

	return c.BindAndValidateResp(&req, func() (*context.GResponse, errors.CommonError) {
		res, signErr := h.surgerySvc.Save(c.Request().Context(), &req, int64(petID), c.GetUserID())
		if signErr != nil {
			return nil, signErr
		}

		return &context.GResponse{
			Status:   http.StatusCreated,
			Response: res,
		}, nil
	})
}

func (h *Handler) GetSurgeriesByPet(c context.Context) error {
	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid pet_id"))
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		res, getErr := h.surgerySvc.GetAllByPet(c.Request().Context(), int64(petID), c.GetUserID())
		if getErr != nil {
			return nil, getErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: surgeries.Surgeries{Surgeries: res},
		}, nil
	})
}

func (h *Handler) GetSurgeryByID(c context.Context) error {
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
		res, getErr := h.surgerySvc.GetById(c.Request().Context(), int64(vaccineID), int64(petID), userID)
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
	var req surgeries.UpdateRequest

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
		res, delErr := h.surgerySvc.UpdateById(c.Request().Context(), int64(vaccineID), &req)
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status:   http.StatusOK,
			Response: res,
		}, nil
	})
}

func (h *Handler) DeleteVetVisitByID(c context.Context) error {
	vaccineID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	petID, err := strconv.Atoi(c.Param("pet_id"))
	if err != nil {
		return errors.NewBadRequestError("invalid parameter", fmt.Errorf("invalid id"))
	}

	return c.NoBindResp(func() (*context.GResponse, errors.CommonError) {
		delErr := h.surgerySvc.DeleteById(c.Request().Context(), int64(vaccineID), int64(petID))
		if delErr != nil {
			return nil, delErr
		}

		return &context.GResponse{
			Status: http.StatusOK,
		}, nil
	})
}
