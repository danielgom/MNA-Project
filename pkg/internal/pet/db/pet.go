package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"MNA-project/pkg/config"
	"MNA-project/pkg/internal/pet/model"
	"MNA-project/pkg/util/errors"
)

type petRepo struct {
	DB config.DBConn
}

// NewPetRepository creates a new pet repository instance.
func NewPetRepository(conn *pgxpool.Pool) Repository {
	return &petRepo{DB: conn}
}

func (r *petRepo) FindByID(ctx context.Context, petID int64, userID int64) (*model.Pet, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM pets WHERE id=$1 AND user_id=$2`, petID, userID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	pets, err := pgx.CollectOneRow(rows, getPetRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("pet not found")
		}
	}

	return pets, nil
}

func (r *petRepo) FindAllByUser(ctx context.Context, userID int64) ([]*model.Pet, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM pets WHERE user_id=$1`, userID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	pets, err := pgx.CollectRows(rows, getPetRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("pets not found")
		}
	}

	return pets, nil
}

func (r *petRepo) Save(ctx context.Context, pet *model.Pet) (*model.Pet, errors.CommonError) {
	row := r.DB.QueryRow(ctx, `INSERT INTO pets(user_id, name, age, breed, birth_date, register_date, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`, pet.UserID, pet.Name, pet.Age, pet.Breed, pet.BirthDate,
		pet.RegisterDate, pet.UpdatedAt)

	saveErr := row.Scan(&pet.ID)
	if saveErr != nil {
		return nil, errors.WrapDatabaseError(saveErr)
	}

	return pet, nil
}

func (r *petRepo) Update(ctx context.Context, pet *model.Pet) errors.CommonError {
	query, parameters := buildUpdateQuery(pet)

	_, err := r.DB.Exec(ctx, query, parameters...)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func (r *petRepo) Delete(ctx context.Context, petID int64, userID int64) errors.CommonError {
	_, err := r.DB.Exec(ctx, `DELETE FROM pets WHERE id=$1 and user_id=$2`, petID, userID)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func getPetRow(row pgx.CollectableRow) (*model.Pet, error) {
	pet, err := pgx.RowToStructByName[model.Pet](row)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

func buildUpdateQuery(pet *model.Pet) (string, []any) {
	nextParam := 2
	parameters := make([]any, 0)

	query := "UPDATE pets SET updated_at=$1"
	parameters = append(parameters, pet.UpdatedAt)

	if pet.Name != "" {
		query += ", name=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, pet.Name)
	}
	if pet.Age != 0 {
		query += ", age=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, pet.Age)
	}
	if pet.Breed != "" {
		query += ", breed=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, pet.Breed)
	}

	parameters = append(parameters, pet.ID)

	return query + " WHERE id=$" + strconv.Itoa(nextParam), parameters
}
