package repository

import (
	"context"
	"fmt"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
)

type (
	OrderDetailQueryRepository struct {
		db pgx.Tx
	}
)

func NewOrderDetailQueryRepository(db pgx.Tx) *OrderDetailQueryRepository {
	return &OrderDetailQueryRepository{db}
}
func (r *OrderDetailQueryRepository) GetOrderDetails(ctx context.Context) ([]model.OrderDetail, error) {

	orderDetails := []model.OrderDetail{}
	query := ` SELECT order_number, product_code,quantity_ordered, price_each, order_line_number FROM orderdetails; `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.OrderDetail
		if err := rows.Scan(
			&a.OrderNumber, &a.ProductCode, &a.QuantityOrdered, &a.PriceEach, &a.OrderLineNumber,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}
		orderDetails = append(orderDetails, a)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", rows.Err())
	}

	return orderDetails, nil
}

func (r *OrderDetailQueryRepository) FindOrderDetails(ctx context.Context, orderNumber int) ([]model.OrderDetail, error) {
	orderDetails := []model.OrderDetail{}

	rows, err := r.db.Query(ctx, `SELECT order_number, product_code,quantity_ordered, price_each, order_line_number FROM orderdetails WHERE order_number=$1`, orderNumber)

	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.OrderDetail
		if err := rows.Scan(
			&a.OrderNumber, &a.ProductCode, &a.QuantityOrdered, &a.PriceEach, &a.OrderLineNumber,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}
		orderDetails = append(orderDetails, a)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", rows.Err())
	}
	return orderDetails, err

}
func (r *OrderDetailQueryRepository) DeleteOrder(ctx context.Context, orderNumber int) error {

	rows, err := r.db.Query(ctx, `DELETE FROM orderdetails WHERE order_number=$1`, orderNumber)

	if err != nil {
		return fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	return err

}
