package repository

import (
	"context"
	"fmt"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
)

type (
	CustomerQueryRepository struct {
		db pgx.Tx
	}
)

func NewCustomerQueryRepository(db pgx.Tx) *CustomerQueryRepository {
	return &CustomerQueryRepository{db}
}

func (r *CustomerQueryRepository) GetCustomers(ctx context.Context) ([]model.Customer, error) {

	customers := []model.Customer{}
	query := `SELECT customer_number, customer_name, contact_last_name, contact_first_name, phone, address_line1,COALESCE(address_line2,''), city , COALESCE(state,''), country, COALESCE(sales_rep_employee_number,0),
	credit_limit FROM customers;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Customer
		if err := rows.Scan(
			&a.CustomerNumber, &a.CustomerName, &a.ContactLastName, &a.ContactFirstName, &a.Phone, &a.AddressLine, &a.AddressLine2, &a.City, &a.State, &a.Country, &a.SalesRepEmpNumber, &a.CreditLimit,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}
		customers = append(customers, a)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", err)
	}

	return customers, nil
}

func (r *CustomerQueryRepository) FindCustomer(ctx context.Context, customerNumber int) (*model.Customer, error) {

	var a model.Customer
	rows, err := r.db.Query(ctx, `SELECT customer_number, customer_name, contact_last_name, contact_first_name, phone, address_line1,COALESCE(address_line2,''), city , COALESCE(state,''), country, COALESCE(sales_rep_employee_number,0),
	credit_limit FROM customers WHERE customer_number=$1`, customerNumber)

	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(
			&a.CustomerNumber, &a.CustomerName, &a.ContactLastName, &a.ContactFirstName, &a.Phone, &a.AddressLine, &a.AddressLine2, &a.City, &a.State, &a.Country, &a.SalesRepEmpNumber, &a.CreditLimit,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}

	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", err)
	}
	return &a, err

}

func (r *CustomerQueryRepository) DeleteCustomer(ctx context.Context, customerNumber int) error {

	rows, err := r.db.Query(ctx, `DELETE FROM customers WHERE customer_number=$1`, customerNumber)

	if err != nil {
		return fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	return err

}
