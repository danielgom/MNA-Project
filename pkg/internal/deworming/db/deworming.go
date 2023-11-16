package db

import (
	"MNA-project/pkg/internal/deworming/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"

	"MNA-project/pkg/config"
	"MNA-project/pkg/util/errors"
)

type dewormingRepo struct {
	DB config.DBConn
}

func NewDewormingRepository(conn *pgxpool.Pool) Repository {
	return &dewormingRepo{DB: conn}
}

func (r *dewormingRepo) FindByID(ctx context.Context, dewormingID int64, petID int64) (*model.Deworming, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM dewormings WHERE id=$1 AND pet_id=$2`, dewormingID, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	vetVisits, err := pgx.CollectOneRow(rows, getDewormingRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("vaccine found")
		}
	}

	return vetVisits, nil
}

func (r *dewormingRepo) FindAllByPet(ctx context.Context, petID int64) ([]*model.Deworming, errors.CommonError) {
	fmt.Println(petID)
	rows, err := r.DB.Query(ctx, `SELECT * FROM dewormings WHERE pet_id=$1`, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	vetVisits, err := pgx.CollectRows(rows, getDewormingRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("pets not found")
		}
	}

	return vetVisits, nil
}

func (r *dewormingRepo) Save(ctx context.Context, deworming *model.Deworming) (*model.Deworming, errors.CommonError) {
	row := r.DB.QueryRow(ctx, `INSERT INTO dewormings(pet_id, vet_name, address, date, next_date, updated_at)
		VALUES($1, $2, $3, $4, $5, $6) RETURNING id`, deworming.PetID, deworming.VetName, deworming.Address, deworming.Date,
		deworming.NextDate, deworming.UpdatedAt)

	saveErr := row.Scan(&deworming.ID)
	if saveErr != nil {
		return nil, errors.WrapDatabaseError(saveErr)
	}

	return deworming, nil
}

func (r *dewormingRepo) Update(ctx context.Context, deworming *model.Deworming) errors.CommonError {
	query, parameters := buildUpdateQuery(deworming)

	_, err := r.DB.Exec(ctx, query, parameters...)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func (r *dewormingRepo) Delete(ctx context.Context, dewormingID int64, userID int64) errors.CommonError {
	_, err := r.DB.Exec(ctx, `DELETE FROM dewormings WHERE id=$1 and pet_id=$2`, dewormingID, userID)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func getDewormingRow(row pgx.CollectableRow) (*model.Deworming, error) {
	vaccine, err := pgx.RowToStructByName[model.Deworming](row)
	if err != nil {
		return nil, err
	}

	return &vaccine, nil
}

func buildUpdateQuery(vac *model.Deworming) (string, []any) {
	nextParam := 2
	parameters := make([]any, 0)

	query := "UPDATE dewormings SET updated_at=$1"
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
