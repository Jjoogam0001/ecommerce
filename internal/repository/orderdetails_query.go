package repository

import (
	"context"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"emperror.dev/errors"
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
		return nil, errors.Wrap(err, "error executing query")
	}

	defer rows.Close()

	for rows.Next() {
		var a model.OrderDetail
		if err := rows.Scan(
			&a.OrderNumber, &a.ProductCode, &a.QuantityOrdered, &a.PriceEach, &a.OrderLineNumber,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning rows")
		}
		orderDetails = append(orderDetails, a)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "error while reading")
	}

	return orderDetails, nil
}
