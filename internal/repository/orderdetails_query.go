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
	OrderDetailQueryRepository struct {
		db pgx.Tx
	}
)

func NewOrderDetailQueryRepository(db pgx.Tx) *OrderDetailQueryRepository {
	return &OrderDetailQueryRepository{db}
}
func (r *OrderDetailQueryRepository) GetOrderDetails(ctx context.Context) ([]model.OrderDetail, error) {
	defer metrics.DBCallSince(time.Now())
	orderDetails := []model.OrderDetail{}
	query := ` SELECT order_number, product_code,quantity_ordered, price_each, order_line_number FROM orderdetails; `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.OrderDetail
		if err := rows.Scan(
			&a.OrderNumber, &a.ProductCode, &a.QuantityOrdered, &a.PriceEach, &a.OrderLineNumber,
		); err != nil {
			return nil, err
		}
		orderDetails = append(orderDetails, a)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return orderDetails, nil
}

func (r *OrderDetailQueryRepository) FindOrderDetails(ctx context.Context, orderNumber int) ([]model.OrderDetail, error) {
	defer metrics.DBCallSince(time.Now())
	orderDetails := []model.OrderDetail{}

	rows, err := r.db.Query(ctx, `SELECT order_number, product_code,quantity_ordered, price_each, order_line_number FROM orderdetails WHERE order_number=$1`, orderNumber)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.OrderDetail
		if err := rows.Scan(
			&a.OrderNumber, &a.ProductCode, &a.QuantityOrdered, &a.PriceEach, &a.OrderLineNumber,
		); err != nil {
			return nil, err
		}
		orderDetails = append(orderDetails, a)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return orderDetails, err

}
func (r *OrderDetailQueryRepository) DeleteOrder(ctx context.Context, orderNumber int) error {
	defer metrics.DBCallSince(time.Now())
	rows, err := r.db.Query(ctx, `DELETE FROM orderdetails WHERE order_number=$1`, orderNumber)

	if err != nil {
		return err
	}

	defer rows.Close()

	return err

}
func (r *OrderDetailQueryRepository) UpdateOrderDetails(ctx context.Context, ordersDetails model.OrderDetail) error {

	sql := `INSERT INTO orderdetails (order_number, product_code, 
                                  quantity_ordered,price_each,order_line_number) values ($1,$2,$3,$4,$5) 
                                  ON CONFLICT (order_number) DO UPDATE SET
                                  order_number=$1,
                                  product_code=$2,
                                  quantity_ordered=$3,
                                  price_each=$4,
                                  order_line_number=$5;`

	log.Infof("Request to update ordersDetails [%v]", ordersDetails)
	_, err := r.db.Exec(ctx, sql,
		ordersDetails.OrderNumber, ordersDetails.ProductCode, ordersDetails.QuantityOrdered,
		ordersDetails.PriceEach, ordersDetails.OrderLineNumber)

	if err != nil {
		return fmt.Errorf("Error while updating ordersDetails", err.Error())
	}
	return nil
}
