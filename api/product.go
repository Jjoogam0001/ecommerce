package api

import (
	"context"
	"net/http"

	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/api/middleware"
	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/internal/repository"
	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/model"
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
