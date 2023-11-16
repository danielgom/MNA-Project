package db

import (
	"MNA-project/pkg/internal/vet_visits/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"

	"MNA-project/pkg/config"
	"MNA-project/pkg/util/errors"
)

type vetVisitRepo struct {
	DB config.DBConn
}

func NewVetVisitRepository(conn *pgxpool.Pool) Repository {
	return &vetVisitRepo{DB: conn}
}

func (r *vetVisitRepo) FindByID(ctx context.Context, vetVisitID int64, petID int64) (*model.VetVisit, errors.CommonError) {
	rows, err := r.DB.Query(ctx, `SELECT * FROM vet_visits WHERE id=$1 AND pet_id=$2`, vetVisitID, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	vetVisits, err := pgx.CollectOneRow(rows, getVisitRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("vaccine found")
		}
	}

	return vetVisits, nil
}

func (r *vetVisitRepo) FindAllByPet(ctx context.Context, petID int64) ([]*model.VetVisit, errors.CommonError) {
	fmt.Println(petID)
	rows, err := r.DB.Query(ctx, `SELECT * FROM vet_visits WHERE pet_id=$1`, petID)
	if err != nil {
		return nil, errors.WrapDatabaseError(err)
	}

	vetVisits, err := pgx.CollectRows(rows, getVisitRow)
	if err != nil {
		if strings.Contains(err.Error(), errors.ErrNotFound) {
			return nil, errors.NewNotFoundError("pets not found")
		}
	}

	return vetVisits, nil
}

func (r *vetVisitRepo) Save(ctx context.Context, vetVisit *model.VetVisit) (*model.VetVisit, errors.CommonError) {
	row := r.DB.QueryRow(ctx, `INSERT INTO vet_visits(pet_id, vet_name, address, reason, comments, date, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id`, vetVisit.PetID, vetVisit.VetName, vetVisit.Address,
		vetVisit.Reason, vetVisit.Comments, vetVisit.Date, vetVisit.UpdatedAt)

	saveErr := row.Scan(&vetVisit.ID)
	if saveErr != nil {
		return nil, errors.WrapDatabaseError(saveErr)
	}

	return vetVisit, nil
}

func (r *vetVisitRepo) Update(ctx context.Context, vetVisit *model.VetVisit) errors.CommonError {
	query, parameters := buildUpdateQuery(vetVisit)

	_, err := r.DB.Exec(ctx, query, parameters...)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func (r *vetVisitRepo) Delete(ctx context.Context, vetVisitID int64, userID int64) errors.CommonError {
	_, err := r.DB.Exec(ctx, `DELETE FROM vet_visits WHERE id=$1 and pet_id=$2`, vetVisitID, userID)
	if err != nil {
		return errors.WrapDatabaseError(err)
	}

	return nil
}

func getVisitRow(row pgx.CollectableRow) (*model.VetVisit, error) {
	vaccine, err := pgx.RowToStructByName[model.VetVisit](row)
	if err != nil {
		return nil, err
	}

	return &vaccine, nil
}

func buildUpdateQuery(vac *model.VetVisit) (string, []any) {
	nextParam := 2
	parameters := make([]any, 0)

	query := "UPDATE vet_visits SET updated_at=$1"
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
	if vac.Reason != "" {
		query += ", reason=$" + strconv.Itoa(nextParam)
		nextParam++
		parameters = append(parameters, vac.Reason)
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
