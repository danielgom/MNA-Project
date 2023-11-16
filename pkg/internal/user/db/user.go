// Package db contains all repositories used by this API.
package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"MNA-project/pkg/config"
	"MNA-project/pkg/internal/user/model"
	"MNA-project/pkg/util/errors"
)

const errNotFound = "no rows in result set"

type userRepo struct {
	DB config.DBConn
}

// NewUserRepository creates a new user repository instance.
func NewUserRepository(conn *pgxpool.Pool) Repository {
	return &userRepo{DB: conn}
}

// FindByEmail finds a user by its email.
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*model.User, errors.CommonError) {
	return r.findUser(ctx, `SELECT * FROM users WHERE email=$1`, email)
}

// FindById finds a user by its id.
func (r *userRepo) FindById(ctx context.Context, id int64) (*model.User, errors.CommonError) {
	return r.findUser(ctx, `SELECT * FROM users WHERE id = $1`, id)
}

// FindAll finds a all users saved .
func (r *userRepo) FindAll(ctx context.Context) ([]*model.User, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM users`)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	users, err := pgx.CollectRows(rows, getUserRow)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	return users, nil
}

// Save persists a user to the DB.
func (r *userRepo) Save(ctx context.Context, user *model.User) (*model.User, errors.CommonError) {
	row := r.DB.QueryRow(ctx, `INSERT INTO users("email", "password", "name", "last_name",
                  "last_login", "address", "created_at", updated_at) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		user.Email, user.Password, user.Name, user.LastName, user.LastLogin, user.Address, user.CreatedAt, user.UpdatedAt)

	saveErr := row.Scan(&user.ID)
	if saveErr != nil {
		return nil, errors.WrapDatabaseError(saveErr)
	}

	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) errors.CommonError {
	query, parameters := buildUpdateQuery(user)

	_, err := r.DB.Exec(ctx, query, parameters...)

	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func (r *userRepo) Delete(ctx context.Context, id int64) errors.CommonError {
	_, err := r.DB.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)

	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func (r *userRepo) findUser(ctx context.Context, query string, args ...any) (*model.User, errors.CommonError) {
	row, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	user, err := pgx.CollectOneRow(row, getUserRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("user not found")
		}

		return nil, errors.WrapDatabaseError(err)
	}

	return user, nil
}

func getUserRow(row pgx.CollectableRow) (*model.User, error) {
	user, err := pgx.RowToStructByName[model.User](row)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func buildUpdateQuery(user *model.User) (string, []any) {
	nextParam := 2
	parameters := make([]any, 0)

	query := "UPDATE users SET updated_at=$1"
	parameters = append(parameters, user.UpdatedAt)

	if user.LastName != "" {
		query += ", last_name=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, user.LastName)
	}
	if user.Name != "" {
		query += ", name=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, user.Name)
	}
	if user.Email != "" {
		query += ", email=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, user.Email)
	}
	if user.Address != "" {
		query += ", address=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, user.Address)
	}

	parameters = append(parameters, user.ID)

	return query + " WHERE id=$" + strconv.Itoa(nextParam), parameters
}
