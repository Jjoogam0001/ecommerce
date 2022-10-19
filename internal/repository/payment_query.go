package repository

import (
	"context"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"emperror.dev/errors"
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

	payments := []model.Payment{}
	query := ` SELECT  customer_number, check_number, payment_date, amount FROM payments; `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Payment
		if err := rows.Scan(
			&a.CustomerNumber, &a.CheckNumber, &a.PaymentDate, &a.Amount,
		); err != nil {
			return nil, errors.Errorf("error scanning rows", err)
		}
		payments = append(payments, a)
	}

	if rows.Err() != nil {
		return nil, errors.Errorf("error while reading", rows.Err())
	}

	return payments, nil
}

func (r *PaymentQueryRepository) FindPayment(ctx context.Context, customerNumber int) ([]model.Payment, error) {
	var pymts []model.Payment

	rows, err := r.db.Query(ctx, `SELECT  customer_number, check_number, payment_date, amount FROM payments WHERE customer_number=$1`, customerNumber)

	if err != nil {
		return nil, errors.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Payment
		if err := rows.Scan(
			&a.CustomerNumber, &a.CheckNumber, &a.PaymentDate, &a.Amount,
		); err != nil {
			return nil, errors.Errorf("error scanning rows", err)
		}

		pymts = append(pymts, a)
	}
	if rows.Err() != nil {
		return nil, errors.Errorf("error while reading", rows.Err())
	}
	return pymts, err

}
func (r *PaymentQueryRepository) DeletePayment(ctx context.Context, customerNumber int) error {

	rows, err := r.db.Query(ctx, `DELETE FROM payments WHERE customer_number=$1`, customerNumber)

	if err != nil {
		return errors.Errorf("error executing query", err)
	}

	defer rows.Close()

	return err

}
