package repository

import (
	"context"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"fmt"
	"time"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/gommon/log"
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
	defer metrics.DBCallSince(time.Now())
	employees := []model.Employee{}
	query := `SELECT employee_number,first_name,last_name,extension,email,office_code, COALESCE(reports_to,0), job_title FROM employees ;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Employee
		if err := rows.Scan(
			&a.EmployeeNumber, &a.FirstName, &a.LastName, &a.Extension, &a.Email, &a.OfficeCode, &a.ReportsTo, &a.Job_Title,
		); err != nil {
			return nil, err
		}
		employees = append(employees, a)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return employees, nil
}
func (r *EmployeeQueryRepository) FindEmployee(ctx context.Context, employeeNumber int) (*model.Employee, error) {
	defer metrics.DBCallSince(time.Now())
	var a model.Employee
	rows, err := r.db.Query(ctx, `SELECT employee_number,first_name,last_name,extension,email,office_code, COALESCE(reports_to,0), job_title FROM employees  WHERE employee_number=$1`, employeeNumber)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(
			&a.EmployeeNumber, &a.FirstName, &a.LastName, &a.Extension, &a.Email, &a.OfficeCode, &a.ReportsTo, &a.Job_Title,
		); err != nil {
			return nil, err
		}

	}
	if rows.Err() != nil {
		return nil, err
	}
	return &a, err

}
func (r *EmployeeQueryRepository) DeleteEmployee(ctx context.Context, employeeNumber int) error {
	defer metrics.DBCallSince(time.Now())
	rows, err := r.db.Query(ctx, `DELETE FROM employees  WHERE employee_number=$1`, employeeNumber)

	if err != nil {
		return err
	}

	defer rows.Close()

	return err

}
func (r *EmployeeQueryRepository) UpdateEmployee(ctx context.Context, employee model.Employee) error {
	defer metrics.DBCallSince(time.Now())
	sql := `INSERT INTO employees (last_name, first_name, 
                                  extension,email,office_code,reports_to,job_title) values ($1,$2,$3,$4,$5,$6,$7) 
                                  ON CONFLICT (email) DO UPDATE SET
                                  last_name=$1,
                                  first_name=$2,
                                  extension=$3,
                                  email=$4,
                                  office_code=$5,
                                  reports_to=$6,
                                  job_title=$7;`

	log.Infof("Request to update employee [%v]", employee)
	_, err := r.db.Exec(ctx, sql,
		employee.LastName, employee.FirstName, employee.Extension,
		employee.Email, employee.OfficeCode, employee.ReportsTo,
		employee.Job_Title)

	if err != nil {
		return fmt.Errorf("Error while updating employee", err.Error())
	}
	return nil
}
