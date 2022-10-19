package repository

import (
	"context"
	"fmt"

	"dev.azure.com/jjoogam/Ecommerce-core/model"

	"github.com/jackc/pgx/v4"
)

type (
	OfficeQueryRepository struct {
		db pgx.Tx
	}
)

func NewOfficeQueryRepository(db pgx.Tx) *OfficeQueryRepository {
	return &OfficeQueryRepository{db}
}
func (r *OfficeQueryRepository) GetOffices(ctx context.Context) ([]model.Office, error) {

	offices := []model.Office{}
	query := ` SELECT office_code, city, phone, address_line1, address_line2, state, country FROM offices; `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Office
		if err := rows.Scan(
			&a.OfficeCode, &a.City, &a.Phone, &a.AddressLine, &a.AddressLine2, &a.State, &a.Country,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}
		offices = append(offices, a)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", err)
	}

	return offices, nil
}

func (r *OfficeQueryRepository) FindOffice(ctx context.Context, officeCode string) (*model.Office, error) {

	var a model.Office
	rows, err := r.db.Query(ctx, `SELECT office_code, city, phone, address_line1, address_line2, state, country FROM offices  WHERE office_code=$1`, officeCode)

	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(
			&a.OfficeCode, &a.City, &a.Phone, &a.AddressLine, &a.AddressLine2, &a.State, &a.Country,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}

	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", rows.Err())
	}
	return &a, err

}

func (r *OfficeQueryRepository) DeleteOffice(ctx context.Context, officeCode string) error {

	rows, err := r.db.Query(ctx, `DELETE  FROM offices  WHERE office_code=$1`, officeCode)

	if err != nil {
		fmt.Errorf("error executing query", err)
	}

	defer rows.Close()
	return err

}
