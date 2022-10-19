package repository

import (
	"context"
	"fmt"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
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
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Employee
		if err := rows.Scan(
			&a.EmployeeNumber, &a.FirstName, &a.LastName, &a.Extension, &a.Email, &a.OfficeCode, &a.ReportsTo, &a.Job_Title,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}
		employees = append(employees, a)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", err)
	}

	return employees, nil
}
func (r *EmployeeQueryRepository) FindEmployee(ctx context.Context, employeeNumber int) (*model.Employee, error) {

	var a model.Employee
	rows, err := r.db.Query(ctx, `SELECT employee_number,first_name,last_name,extension,email,office_code, COALESCE(reports_to,0), job_title FROM employees  WHERE employee_number=$1`, employeeNumber)

	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(
			&a.EmployeeNumber, &a.FirstName, &a.LastName, &a.Extension, &a.Email, &a.OfficeCode, &a.ReportsTo, &a.Job_Title,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}

	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", err)
	}
	return &a, err

}
func (r *EmployeeQueryRepository) DeleteEmployee(ctx context.Context, employeeNumber int) error {

	rows, err := r.db.Query(ctx, `DELETE FROM employees  WHERE employee_number=$1`, employeeNumber)

	if err != nil {
		return fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	return err

}
