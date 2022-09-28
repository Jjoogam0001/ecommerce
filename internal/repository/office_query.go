package repository

import (
	"context"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"emperror.dev/errors"
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
		return nil, errors.Wrap(err, "error executing query")
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Office
		if err := rows.Scan(
			&a.OfficeCode, &a.City, &a.Phone, &a.AddressLine, &a.AddressLine2, &a.State, &a.Country,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning rows")
		}
		offices = append(offices, a)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "error while reading")
	}

	return offices, nil
}
