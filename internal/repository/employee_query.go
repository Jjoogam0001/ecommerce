package repository

import (
	"context"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"emperror.dev/errors"
	"github.com/jackc/pgx/v4"
)

type (
	EmployeeQueryRepository struct {
		db pgx.Tx
	}
)

func NewEmployeeQueryRepository(db pgx.Tx) *EmployeeQueryRepository {
	return &EmployeeQueryRepository{db}
}

func (r *EmployeeQueryRepository) GetEmployees(ctx context.Context) ([]model.Employee, error) {

	employees := []model.Employee{}
	query := `SELECT employee_number,first_name,last_name,extension,email,office_code, COALESCE(reports_to,0), job_title FROM employees ;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Employee
		if err := rows.Scan(
			&a.EmployeeNumber, &a.FirstName, &a.LastName, &a.Extension, &a.Email, &a.OfficeCode, &a.ReportsTo, &a.Job_Title,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning rows")
		}
		employees = append(employees, a)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "error while reading")
	}

	return employees, nil
}
