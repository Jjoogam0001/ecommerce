package repository

import (
	"context"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"emperror.dev/errors"
	"github.com/jackc/pgx/v4"
)

type (
	ProductQueryRepository struct {
		db pgx.Tx
	}
)

func NewProductQueryRepository(db pgx.Tx) *ProductQueryRepository {
	return &ProductQueryRepository{db}
}

func (r *ProductQueryRepository) GetProducts(ctx context.Context) ([]model.Product, error) {

	products := []model.Product{}
	query := `SELECT  product_code, product_name, product_line,product_scale, product_vendor, product_description, quantity_in_stock, buy_price, msrp FROM products;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "error executing query")
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Product
		if err := rows.Scan(
			&a.ProductCode, &a.ProductName, &a.ProductLine, &a.ProductScale, &a.ProductVendor, &a.ProductDescription, &a.QuantityInStock, &a.BuyPrice, &a.Msrp,
		); err != nil {
			return nil, errors.Wrap(err, "error scanning rows")
		}
		products = append(products, a)
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "error while reading")
	}

	return products, nil
}
