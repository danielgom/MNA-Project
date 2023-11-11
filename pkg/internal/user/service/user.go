// Package service contains all core logic for this API.
package service

import (
	"context"
	"net/http"
	"net/mail"
	"time"

	"MNA-project/pkg/internal/user"
	"MNA-project/pkg/internal/user/db"
	"MNA-project/pkg/internal/user/model"
	"MNA-project/pkg/security"
	"MNA-project/pkg/util/errors"
)

type userSvc struct {
	userDB db.Repository
}

// NewUserService returns a new instance of user service.
func NewUserService(uR db.Repository) UserService {
	return &userSvc{userDB: uR}
}

// SignUp executes core logic in order to save the user and generate its verification token for the first time.
func (u *userSvc) SignUp(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse,
	errors.CommonError) {
	pass, err := security.Hash(req.Password)
	if err != nil {
		return nil, errors.NewRestError("password not encrypted",
			http.StatusInternalServerError, "Internal server error", err)
	}

	currentTime := time.Now().Local()

	newUser := new(model.User)

	newUser.Email = req.Email
	newUser.Name = req.Name
	newUser.LastName = req.LastName
	newUser.Password = pass
	newUser.CreatedAt = &currentTime
	newUser.UpdatedAt = &currentTime

	newUser, saveErr := u.userDB.Save(ctx, newUser)
	if saveErr != nil {
		return nil, saveErr
	}

	return user.BuildRegisterResponse(newUser), nil
}

func (u *userSvc) GetById(ctx context.Context, id int64) (*user.Response, errors.CommonError) {
	currentUser, commonError := u.userDB.FindById(ctx, id)
	if commonError != nil {
		return nil, commonError
	}

	return &user.Response{
		ID:        currentUser.ID,
		Email:     currentUser.Email,
		Name:      currentUser.Name,
		LastName:  currentUser.LastName,
		LastLogin: currentUser.LastLogin,
		CreatedAt: currentUser.CreatedAt,
		UpdatedAt: currentUser.UpdatedAt,
	}, nil
}

func (u *userSvc) GetAll(ctx context.Context) ([]*user.Response, errors.CommonError) {
	users, commonError := u.userDB.FindAll(ctx)
	if commonError != nil {
		return nil, commonError
	}

	usersResponse := make([]*user.Response, 0, len(users))
	for _, us := range users {
		userRes := &user.Response{
			ID:        us.ID,
			Email:     us.Email,
			Name:      us.Name,
			LastName:  us.LastName,
			LastLogin: us.LastLogin,
			CreatedAt: us.CreatedAt,
			UpdatedAt: us.UpdatedAt,
		}
		usersResponse = append(usersResponse, userRes)
	}
	return usersResponse, nil
}

func (u *userSvc) UpdateById(ctx context.Context, id int64, request *user.UpdateRequest) (*user.Response, errors.CommonError) {
	updateTime := time.Now().Local()
	dbErr := u.userDB.Update(ctx, &model.User{
		ID:        id,
		Email:     request.Email,
		Name:      request.Name,
		LastName:  request.LastName,
		UpdatedAt: &updateTime,
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return u.GetById(ctx, id)
}

func (u *userSvc) Login(ctx context.Context, request *user.LoginRequest) (*user.LoginResponse, errors.CommonError) {
	_, err := mail.ParseAddress(request.Email)
	if err != nil {
		return nil, errors.NewBadRequestError("Invalid email address", err)
	}

	us, commonError := u.userDB.FindByEmail(ctx, request.Email)
	if commonError != nil {
		return nil, commonError
	}

	validPass := security.CheckHash(request.Password, us.Password)
	if !validPass {
		return nil, errors.NewUnauthorisedError("invalid password")
	}

	JWT, expDate, err := security.GenerateTokenWithExp(us.ID)
	if err != nil {
		return nil, errors.NewInternalServerError("internal error", err)
	}

	return user.BuildLoginResponse(us, expDate, JWT), nil
}

func (u *userSvc) DeleteById(ctx context.Context, id int64) errors.CommonError {
	return u.userDB.Delete(ctx, id)
}
