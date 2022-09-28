package repository

import (
	"context"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"emperror.dev/errors"
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
		return nil, errors.Wrap(err, "error executing query")
	}

	defer rows.Close()

	for rows.Next() {
		var a model.ProductLine
		if err := rows.Scan(
			&a.ProductLine, &a.TextDescription, &a.Image,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning rows")
		}
		productLines = append(productLines, a)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "error while reading")
	}

	return productLines, nil
}
