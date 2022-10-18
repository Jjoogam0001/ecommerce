package api

import (
	"context"
	"net/http"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/repository"
	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
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
	productGroup.GET("/get_products", a.getorders)
	productGroup.GET("/get_product", a.findProduct)
	productGroup.DELETE("/delete_product", a.deleteProduct)

}

// @Summary Retrieve all Payments
// @Descript Products
// @Produce json
// @Tags Products
// @Router /products/get_products [get]
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrValidation
func (a *ProductController) getorders(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetProducts(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve a single Product
// @Description Fetch a single Product
// @Tags Products
// @Produce json
// @Router /products/get_product [get]
// @Param product_code query string true "product_code mandatory"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrValidation
func (a *ProductController) findProduct(c echo.Context) error {
	cuid, err := a.decodeProduct(c)
	if err != nil {
		return errors.Wrap(err, "unable to decode")
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindProduct(ctx, *cuid)
	if err != nil {
		return errors.Wrap(err, "cant find product")
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary Deletes a single Product
// @Description Deletes a single Product
// @Tags Products
// @Produce json
// @Router /products/delete_product [delete]
// @Param product_code query string true "product_code mandatory"
// @Success 200 {object} model.Product
// @Failure 400 {object} model.ErrValidation
func (a *ProductController) deleteProduct(c echo.Context) error {
	cuid, err := a.decodeProduct(c)
	if err != nil {
		return errors.Wrap(err, "unable to decode")
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindProduct(ctx, *cuid)
	if err != nil {
		return errors.Wrap(err, "cant find product")
	}
	err = r.DeleteProduct(ctx, *cuid)

	if err != nil {
		return errors.Wrap(err, "cant delete product")
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
