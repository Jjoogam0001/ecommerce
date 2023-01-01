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
	PaymentQueryRepository struct {
		db pgx.Tx
	}
)

func NewPaymentQueryRepository(db pgx.Tx) *PaymentQueryRepository {
	return &PaymentQueryRepository{db}
}
func (r *PaymentQueryRepository) GetPayments(ctx context.Context) ([]model.Payment, error) {
	defer metrics.DBCallSince(time.Now())
	payments := []model.Payment{}
	query := ` SELECT  customer_number, check_number, payment_date, amount FROM payments; `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Payment
		if err := rows.Scan(
			&a.CustomerNumber, &a.CheckNumber, &a.PaymentDate, &a.Amount,
		); err != nil {
			return nil, err
		}
		payments = append(payments, a)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return payments, nil
}

func (r *PaymentQueryRepository) FindPayment(ctx context.Context, customerNumber int) ([]model.Payment, error) {
	defer metrics.DBCallSince(time.Now())
	var pymts []model.Payment

	rows, err := r.db.Query(ctx, `SELECT  customer_number, check_number, payment_date, amount FROM payments WHERE customer_number=$1`, customerNumber)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Payment
		if err := rows.Scan(
			&a.CustomerNumber, &a.CheckNumber, &a.PaymentDate, &a.Amount,
		); err != nil {
			return nil, err
		}

		pymts = append(pymts, a)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return pymts, err

}
func (r *PaymentQueryRepository) DeletePayment(ctx context.Context, customerNumber int) error {
	defer metrics.DBCallSince(time.Now())
	rows, err := r.db.Query(ctx, `DELETE FROM payments WHERE customer_number=$1`, customerNumber)

	if err != nil {
		return err
	}

	defer rows.Close()

	return err

}
func (r *PaymentQueryRepository) UpdatePayment(ctx context.Context, payment model.Payment) error {
	defer metrics.DBCallSince(time.Now())
	sql := `INSERT INTO orderdetails (check_number, 
                                  payment_date,amount) values ($1,$2,$3,) 
                                  ON CONFLICT (customer_number) DO UPDATE SET
                                  check_number=$1,
                                  payment_date=$2,
                                  amount=$3;`

	log.Infof("Request to update payment [%v]", payment)
	_, err := r.db.Exec(ctx, sql,
		payment.CheckNumber, payment.PaymentDate,
		payment.Amount)

	if err != nil {
		return fmt.Errorf("Error while updating payment", err.Error())
	}
	return nil
}
