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
	ProductLineController struct {
		queryRepositoryFactory ProductLineQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	ProductLineQueryRepository interface {
		GetProductLines(ctx context.Context) ([]model.ProductLine, error)
	}

	ProductLineQueryRepositoryFactory func(pgx.Tx) ProductLineQueryRepository
)

func NewProductLineController(txProvider middleware.TxProvider) *ProductLineController {
	return &ProductLineController{
		queryRepositoryFactory: defaultProductLineQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultProductLineQueryRepositoryFactory(tx pgx.Tx) ProductLineQueryRepository {
	return repository.NewProductLineQueryRepository(tx)
}

func (a *ProductLineController) WithQueryRepository(f ProductLineQueryRepositoryFactory) *ProductLineController {
	a.queryRepositoryFactory = f
	return a
}

func (a *ProductLineController) RegisterRoutes(e *echo.Echo) {
	productLineGroup := e.Group("/productLines", middleware.Transaction(a.txProvider))
	productLineGroup.GET("/get_productLines", a.GetProductLines)

}

// @Summary Retrieve all ProductLines
// @Descript ProductLines
// @Produce json
// @Tags ProductLines
// @Router /productLines/get_productLines [get]
// @Success 200 {object} model.ProductLine
// @Failure 400 {object} model.ErrValidation
func (a *ProductLineController) GetProductLines(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	productLines, err := r.GetProductLines(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}

	return c.JSON(http.StatusOK, productLines)
}
