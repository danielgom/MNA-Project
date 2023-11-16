package db

import (
	"MNA-project/pkg/internal/surgeries/model"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"

	"MNA-project/pkg/config"
	"MNA-project/pkg/util/errors"
)

type surgeryRepo struct {
	DB config.DBConn
}

func NewSurgeryRepository(conn *pgxpool.Pool) Repository {
	return &surgeryRepo{DB: conn}
}

func (r *surgeryRepo) FindByID(ctx context.Context, surgeryID int64, petID int64) (*model.Surgery, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM surgeries WHERE id=$1 AND pet_id=$2`, surgeryID, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	surgery, err := pgx.CollectOneRow(rows, getSurgeryRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("vaccine found")
		}
	}

	return surgery, nil
}

func (r *surgeryRepo) FindAllByPet(ctx context.Context, petID int64) ([]*model.Surgery, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM surgeries WHERE pet_id=$1`, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	surgeries, err := pgx.CollectRows(rows, getSurgeryRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("pets not found")
		}
	}

	return surgeries, nil
}

func (r *surgeryRepo) Save(ctx context.Context, surgery *model.Surgery) (*model.Surgery, errors.CommonError) {
	row := r.DB.QueryRow(ctx, `INSERT INTO surgeries(pet_id, vet_name, address, name, comments, date, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`, surgery.PetID, surgery.VetName, surgery.Address,
		surgery.Name, surgery.Comments, surgery.Date, surgery.UpdatedAt)

	saveErr := row.Scan(&surgery.ID)
	if saveErr != nil {
		return nil, errors.WrapDatabaseError(saveErr)
	}

	return surgery, nil
}

func (r *surgeryRepo) Update(ctx context.Context, surgery *model.Surgery) errors.CommonError {
	query, parameters := buildUpdateQuery(surgery)

	_, err := r.DB.Exec(ctx, query, parameters...)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func (r *surgeryRepo) Delete(ctx context.Context, surgeryID int64, userID int64) errors.CommonError {
	_, err := r.DB.Exec(ctx, `DELETE FROM surgeries WHERE id=$1 and pet_id=$2`, surgeryID, userID)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func getSurgeryRow(row pgx.CollectableRow) (*model.Surgery, error) {
	vaccine, err := pgx.RowToStructByName[model.Surgery](row)
	if err != nil {
		return nil, err
	}

	return &vaccine, nil
}

func buildUpdateQuery(vac *model.Surgery) (string, []any) {
	nextParam := 2
	parameters := make([]any, 0)

	query := "UPDATE surgeries SET updated_at=$1"
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
	if vac.Name != "" {
		query += ", name=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.Name)
	}
	if vac.Comments != "" {
		query += ", comments=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.Comments)
	}
	if vac.Date != nil {
		query += ", date=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.Date)
	}

	parameters = append(parameters, vac.ID)

	return query + " WHERE id=$" + strconv.Itoa(nextParam), parameters
}
