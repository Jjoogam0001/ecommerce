package repository

import (
	"context"
	"fmt"

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

	productLines := []model.ProductLine{}
	query := `SELECT  product_line, text_description , image FROM product_lines;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query", err)
	}

	defer rows.Close()

	for rows.Next() {
		var a model.ProductLine
		if err := rows.Scan(
			&a.ProductLine, &a.TextDescription, &a.Image,
		); err != nil {
			return nil, fmt.Errorf("error scanning rows", err)
		}
		productLines = append(productLines, a)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error while reading", rows.Err())
	}

	return productLines, nil
}
