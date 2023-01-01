package repository

import (
	"context"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"time"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
)

type (
	ProductLineQueryRepository struct {
		db pgx.Tx
	}
)

func NewProductLineQueryRepository(db pgx.Tx) *ProductLineQueryRepository {
	return &ProductLineQueryRepository{db}
}

func (r *ProductLineQueryRepository) GetProductLines(ctx context.Context) ([]model.ProductLine, error) {
	defer metrics.DBCallSince(time.Now())
	productLines := []model.ProductLine{}
	query := `SELECT  product_line, text_description , image FROM product_lines;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.ProductLine
		if err := rows.Scan(
			&a.ProductLine, &a.TextDescription, &a.Image,
		); err != nil {
			return nil, err
		}
		productLines = append(productLines, a)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return productLines, nil
}
