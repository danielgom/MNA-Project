package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"

	"MNA-project/pkg/config"
	"MNA-project/pkg/internal/vaccine/model"
	"MNA-project/pkg/util/errors"
)

type vaccineRepo struct {
	DB config.DBConn
}

func NewVaccineRepository(conn *pgxpool.Pool) Repository {
	return &vaccineRepo{DB: conn}
}

func (r *vaccineRepo) FindByID(ctx context.Context, vaccineID int64, petID int64) (*model.Vaccine, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM vaccines WHERE id=$1 AND pet_id=$2`, vaccineID, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	pets, err := pgx.CollectOneRow(rows, getVaccineRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("vaccine found")
		}
	}

	return pets, nil
}

func (r *vaccineRepo) FindAllByPet(ctx context.Context, petID int64) ([]*model.Vaccine, errors.CommonError) {
	fmt.Println(petID)
	rows, err := r.DB.Query(ctx, `SELECT * FROM vaccines WHERE pet_id=$1`, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	pets, err := pgx.CollectRows(rows, getVaccineRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("pets not found")
		}
	}

	return pets, nil
}

func (r *vaccineRepo) Save(ctx context.Context, vac *model.Vaccine) (*model.Vaccine, errors.CommonError) {
	row := r.DB.QueryRow(ctx, `INSERT INTO vaccines(pet_id, type, vet_name, address, date, next_date, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`, vac.PetID, vac.Type, vac.VetName, vac.Address, vac.Date,
		vac.NextDate, vac.UpdatedAt)

	saveErr := row.Scan(&vac.ID)
	if saveErr != nil {
		return nil, errors.WrapDatabaseError(saveErr)
	}

	return vac, nil
}

func (r *vaccineRepo) Update(ctx context.Context, vac *model.Vaccine) errors.CommonError {
	query, parameters := buildUpdateQuery(vac)

	_, err := r.DB.Exec(ctx, query, parameters...)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func (r *vaccineRepo) Delete(ctx context.Context, petID int64, userID int64) errors.CommonError {
	_, err := r.DB.Exec(ctx, `DELETE FROM vaccines WHERE id=$1 and pet_id=$2`, petID, userID)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func getVaccineRow(row pgx.CollectableRow) (*model.Vaccine, error) {
	vaccine, err := pgx.RowToStructByName[model.Vaccine](row)
	if err != nil {
		return nil, err
	}

	return &vaccine, nil
}

func buildUpdateQuery(vac *model.Vaccine) (string, []any) {
	nextParam := 2
	parameters := make([]any, 0)

	query := "UPDATE vaccines SET updated_at=$1"
	parameters = append(parameters, vac.UpdatedAt)

	if vac.Address != "" {
		query += ", address=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.Address)
	}
	if vac.VetName != "" {
		query += ", vet_name=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.VetName)
	}
	if vac.Type != "" {
		query += ", type=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.Type)
	}
	if vac.Date != nil {
		query += ", date=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.Date)
	}
	if vac.NextDate != nil {
		query += ", next_date=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.NextDate)
	}

	parameters = append(parameters, vac.ID)

	return query + " WHERE id=$" + strconv.Itoa(nextParam), parameters
}
