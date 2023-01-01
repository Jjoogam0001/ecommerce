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
	ProductQueryRepository struct {
		db pgx.Tx
	}
)

func NewProductQueryRepository(db pgx.Tx) *ProductQueryRepository {
	return &ProductQueryRepository{db}
}

func (r *ProductQueryRepository) GetProducts(ctx context.Context) ([]model.Product, error) {
	defer metrics.DBCallSince(time.Now())
	products := []model.Product{}
	query := `SELECT  product_code, product_name, product_line,product_scale, product_vendor, product_description, quantity_in_stock, buy_price, msrp FROM products;`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a model.Product
		if err := rows.Scan(
			&a.ProductCode, &a.ProductName, &a.ProductLine, &a.ProductScale, &a.ProductVendor, &a.ProductDescription, &a.QuantityInStock, &a.BuyPrice, &a.Msrp,
		); err != nil {
			return nil, err
		}
		products = append(products, a)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return products, nil
}
func (r *ProductQueryRepository) FindProduct(ctx context.Context, productCode string) (*model.Product, error) {
	defer metrics.DBCallSince(time.Now())
	var a model.Product
	rows, err := r.db.Query(ctx, `SELECT  product_code, product_name, product_line,product_scale, product_vendor, product_description, quantity_in_stock, buy_price, msrp FROM products WHERE product_code=$1`, productCode)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		if err := rows.Scan(
			&a.ProductCode, &a.ProductName, &a.ProductLine, &a.ProductScale, &a.ProductVendor, &a.ProductDescription, &a.QuantityInStock, &a.BuyPrice, &a.Msrp,
		); err != nil {
			return nil, err
		}

	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &a, err

}

func (r *ProductQueryRepository) DeleteProduct(ctx context.Context, productCode string) error {
	defer metrics.DBCallSince(time.Now())
	rows, err := r.db.Query(ctx, `DELETE FROM products WHERE product_code=$1`, productCode)

	if err != nil {
		return err
	}

	defer rows.Close()

	return err

}
func (r *ProductQueryRepository) UpdateProduct(ctx context.Context, product model.Product) error {
	defer metrics.DBCallSince(time.Now())
	sql := `INSERT INTO products (product_name, product_line, 
                                  product_scale,product_vendor,product_description,quantity_in_stock,buy_price,msrp) values ($1,$2,$3,$4,$5,$6,$7,$8) 
                                  ON CONFLICT (product_code) DO UPDATE SET
                                  product_name=$1,
                                  product_line=$2,
                                  product_scale=$3,
                                  product_vendor=$4,
                                  product_description=$5,
                                  quantity_in_stock=$6,
                                  buy_price=$7,
                                  msrp=$8;`

	log.Infof("Request to update product [%v]", product)
	_, err := r.db.Exec(ctx, sql,
		product.ProductName, product.ProductLine, product.ProductScale,
		product.ProductVendor, product.ProductDescription, product.QuantityInStock,
		product.BuyPrice, product.Msrp)

	if err != nil {
		return fmt.Errorf("Error while updating product", err.Error())
	}
	return nil
}
