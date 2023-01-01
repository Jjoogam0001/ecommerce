package repository

import (
	"context"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"fmt"
	"github.com/labstack/gommon/log"
	"time"

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
	defer metrics.DBCallSince(time.Now())
	offices := []model.Office{}
	query := ` SELECT office_code, city, phone, address_line1, address_line2, state, country FROM offices; `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Office
		if err := rows.Scan(
			&a.OfficeCode, &a.City, &a.Phone, &a.AddressLine, &a.AddressLine2, &a.State, &a.Country,
		); err != nil {
			return nil, err
		}
		offices = append(offices, a)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return offices, nil
}

func (r *OfficeQueryRepository) FindOffice(ctx context.Context, officeCode string) (*model.Office, error) {
	defer metrics.DBCallSince(time.Now())
	var a model.Office
	rows, err := r.db.Query(ctx, `SELECT office_code, city, phone, address_line1, address_line2, state, country FROM offices  WHERE office_code=$1`, officeCode)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(
			&a.OfficeCode, &a.City, &a.Phone, &a.AddressLine, &a.AddressLine2, &a.State, &a.Country,
		); err != nil {
			return nil, err
		}

	}
	if rows.Err() != nil {
		return nil, err
	}
	return &a, err

}

func (r *OfficeQueryRepository) DeleteOffice(ctx context.Context, officeCode string) error {
	defer metrics.DBCallSince(time.Now())
	rows, err := r.db.Query(ctx, `DELETE  FROM offices  WHERE office_code=$1`, officeCode)

	if err != nil {
		return err
	}

	defer rows.Close()
	return err

}
func (r *OfficeQueryRepository) UpdateOffice(ctx context.Context, office model.Office) error {
	defer metrics.DBCallSince(time.Now())
	sql := `INSERT INTO offices (office_code, city, 
                                  phone,address_line1,address_line2,state,country,postal_code, territory) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) 
                                  ON CONFLICT (office_code) DO UPDATE SET
                                  office_code=$1,
                                  city=$2,
                                  phone=$3,
                                  address_line1=$4,
                                  address_line2=$5,
                                  state=$6,
                                  country=$7,
                                   postal_code=$8,
                                   territory=$9;`

	log.Infof("Request to update office [%v]", office)
	_, err := r.db.Exec(ctx, sql,
		office.OfficeCode, office.City, office.Phone,
		office.AddressLine, office.AddressLine2, office.State,
		office.Country, office.PostalCode, office.Territory)

	if err != nil {
		return fmt.Errorf("Error while updating office", err.Error())
	}
	return nil
}
