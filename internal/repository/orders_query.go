package repository

import (
	"context"
	"fmt"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
)

type (
	OrderQueryRepository struct {
		db pgx.Tx
	}
)

func NewOrderQueryRepository(db pgx.Tx) *OrderQueryRepository {
	return &OrderQueryRepository{db}
}
func (r *OrderQueryRepository) Getorders(ctx context.Context) ([]model.Order, error) {

	orders := []model.Order{}
	query := ` SELECT order_number,order_date,required_date,shipped_date,status,customer_number,comments from orders where comments is not NULL; `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Order
		if err := rows.Scan(
			&a.OrderNumber, &a.OrderDate, &a.RequiredDate, &a.ShippedDate, &a.Status, &a.Customer_number, &a.Comments,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}
		orders = append(orders, a)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", rows.Err())
	}

	return orders, nil
}
