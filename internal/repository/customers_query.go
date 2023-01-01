package repository

import (
	"context"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/gommon/log"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
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
	defer metrics.DBCallSince(time.Now())
	customers := []model.Customer{}
	query := `SELECT email, customer_number, customer_name, contact_last_name, contact_first_name, phone,
       address_line1,COALESCE(address_line2,''), city , COALESCE(state,''), country, 
       COALESCE(sales_rep_employee_number,0),credit_limit FROM customers;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Customer
		if err := rows.Scan(
			&a.Email, &a.Customer_number, &a.CustomerName, &a.ContactLastName, &a.ContactFirstName, &a.Phone, &a.AddressLine, &a.AddressLine2,
			&a.City, &a.State, &a.Country, &a.SalesRepEmpNumber, &a.CreditLimit,
		); err != nil {
			return nil, err
		}
		customers = append(customers, a)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerQueryRepository) FindCustomer(ctx context.Context, email string) (*model.Customer, error) {
	defer metrics.DBCallSince(time.Now())
	var a model.Customer
	rows, err := r.db.Query(ctx, `SELECT email,customer_name, contact_last_name, contact_first_name, phone,
       address_line1,COALESCE(address_line2,''), city , COALESCE(state,''), country, COALESCE(sales_rep_employee_number,0),
	credit_limit FROM customers WHERE email=$1`, email)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(
			&a.Email, &a.CustomerName, &a.ContactLastName, &a.ContactFirstName, &a.Phone, &a.AddressLine, &a.AddressLine2,
			&a.City, &a.State, &a.Country, &a.SalesRepEmpNumber, &a.CreditLimit,
		); err != nil {
			return nil, err
		}

	}
	if rows.Err() != nil {
		return nil, err
	}
	return &a, err

}

func (r *CustomerQueryRepository) DeleteCustomer(ctx context.Context, email string) error {
	defer metrics.DBCallSince(time.Now())
	rows, err := r.db.Query(ctx, `DELETE FROM customers WHERE email=$1`, email)

	if err != nil {
		return err
	}

	defer rows.Close()

	return err

}

func (r *CustomerQueryRepository) UpdateCustomer(ctx context.Context, customer model.Customer) error {
	defer metrics.DBCallSince(time.Now())
	sql := `INSERT INTO customers (customer_name, contact_last_name, contact_first_name,phone, 
                                  address_line1,address_line2,city,state,postal_code,
                                  country,sales_rep_employee_number,email,credit_limit) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,
                                                                                          $11,$12,$13) 
                                  ON CONFLICT (email) DO UPDATE SET
                                  customer_name=$1, 
                                  contact_last_name=$2,
                                  contact_first_name=$3,
                                  phone=$4,
                                  address_line1=$5,
                                  address_line2=$6,
                                  city=$7,
                                  state=$8,
                                  postal_code=$9,
                                  country=$10,
                                  sales_rep_employee_number=$11,
                                  email = $12,
                                  credit_limit=$13;`

	log.Infof("Request to update customer [%v]", customer)
	_, err := r.db.Exec(ctx, sql,
		customer.CustomerName, customer.ContactLastName, customer.ContactFirstName, customer.Phone,
		customer.AddressLine, customer.AddressLine2, customer.City,
		customer.State, customer.PostalCode, customer.Country, customer.SalesRepEmpNumber, customer.Email,
		customer.CreditLimit)

	if err != nil {
		return fmt.Errorf("Error while updating customer", err.Error())
	}
	return nil
}
