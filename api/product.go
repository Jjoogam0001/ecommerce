package api

import (
	"context"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/repository"
	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type (
	ProductController struct {
		queryRepositoryFactory ProductQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	ProductQueryRepository interface {
		GetProducts(ctx context.Context) ([]model.Product, error)
		FindProduct(ctx context.Context, productCode string) (*model.Product, error)
		DeleteProduct(ctx context.Context, productCode string) error
		UpdateProduct(ctx context.Context, product model.Product) error
	}

	ProductQueryRepositoryFactory func(pgx.Tx) ProductQueryRepository
)

func NewProductController(txProvider middleware.TxProvider) *ProductController {
	return &ProductController{
		queryRepositoryFactory: defaultProductQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultProductQueryRepositoryFactory(tx pgx.Tx) ProductQueryRepository {
	return repository.NewProductQueryRepository(tx)
}

func (a *ProductController) WithQueryRepository(f ProductQueryRepositoryFactory) *ProductController {
	a.queryRepositoryFactory = f
	return a
}

func (a *ProductController) RegisterRoutes(e *echo.Echo) {
	productGroup := e.Group("/products", middleware.Transaction(a.txProvider))
	productGroup.GET("/products", a.getorders)
	productGroup.GET("/get", a.findProduct)
	productGroup.DELETE("/delete", a.deleteProduct)
	productGroup.PATCH("/update", a.UpdateProduct)

}

// @Summary Retrieve all Payments
// @Descript Products
// @Produce json
// @Tags Products
// @Router /products/products [get]
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrValidation
func (a *ProductController) getorders(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetProducts(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve a single Product
// @Description Fetch a single Product
// @Tags Products
// @Produce json
// @Router /products/get [get]
// @Param product_code query string true "product_code mandatory"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrValidation
func (a *ProductController) findProduct(c echo.Context) error {
	cuid, err := a.decodeProduct(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindProduct(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary Deletes a single Product
// @Description Deletes a single Product
// @Tags Products
// @Produce json
// @Router /products/delete [delete]
// @Param product_code query string true "product_code mandatory"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrValidation
func (a *ProductController) deleteProduct(c echo.Context) error {
	cuid, err := a.decodeProduct(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindProduct(ctx, *cuid)
	if err != nil {
		return err
	}
	err = r.DeleteProduct(ctx, *cuid)

	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.ProductResponse{
		Product: *customer,
		Status:  "Deleted",
	})

}

func (h *ProductController) decodeProduct(c echo.Context) (*string, error) {
	customerNumber := c.QueryParam("product_code")
	if customerNumber == "" && len(customerNumber) == 0 {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "product_code", Reason: "Missing key product_code ."}}}
	}

	return &customerNumber, nil
}

// @Summary Updates a  Product
// @Description updates a single Product
// @Tags Products
// @Produce json
// @Router /products/update [patch]
// @Param product body model.Product true "quantity_in_stock"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrValidation
func (a *ProductController) UpdateProduct(c echo.Context) error {
	ctx := context.Background()
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		metrics.DBErrorInc()
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)
	requestbody, err := ioutil.ReadAll(c.Request().Body)

	var request model.Product
	err = json.Unmarshal(requestbody, &request)
	if err != nil {
		return err
	}

	err = r.UpdateProduct(ctx, request)
	if err != nil {
		return fmt.Errorf("Error in the request", err.Error())
	}
	return c.String(http.StatusOK, "Product updated successfully")

}
