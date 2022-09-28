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
	OfficeController struct {
		queryRepositoryFactory OfficeQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	OfficeQueryRepository interface {
		GetOffices(ctx context.Context) ([]model.Office, error)
	}

	OfficeQueryRepositoryFactory func(pgx.Tx) OfficeQueryRepository
)

func NewOfficeController(txProvider middleware.TxProvider) *OfficeController {
	return &OfficeController{
		queryRepositoryFactory: defaultOfficeQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultOfficeQueryRepositoryFactory(tx pgx.Tx) OfficeQueryRepository {
	return repository.NewOfficeQueryRepository(tx)
}

func (a *OfficeController) WithQueryRepository(f OfficeQueryRepositoryFactory) *OfficeController {
	a.queryRepositoryFactory = f
	return a
}

func (a *OfficeController) RegisterRoutes(e *echo.Echo) {
	officeGroup := e.Group("/offices", middleware.Transaction(a.txProvider))
	officeGroup.GET("/get_offices", a.GetOffices)

}

// @Summary Retrieve all Offices
// @Description Gets all Offices
// @Tags 	Offices
// @Produce json
// @Router /offices/get_offices [get]
// @Success 200 {object} model.Office
// @Failure 400 {object} model.ErrValidation
func (a *OfficeController) GetOffices(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetOffices(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}

	return c.JSON(http.StatusOK, orders)
}
