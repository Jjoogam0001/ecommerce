package api

import (
	"context"
	"net/http"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/repository"
	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type (
	StudentController struct {
		queryRepositoryFactory StudentQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	StudentQueryRepository interface {
		Getorders(ctx context.Context) ([]model.Order, error)
	}

	StudentQueryRepositoryFactory func(pgx.Tx) StudentQueryRepository
)

func NewStudentController(txProvider middleware.TxProvider) *StudentController {
	return &StudentController{
		queryRepositoryFactory: defaultStudentQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultStudentQueryRepositoryFactory(tx pgx.Tx) StudentQueryRepository {
	return repository.NewOrderQueryRepository(tx)
}

func (a *StudentController) WithQueryRepository(f StudentQueryRepositoryFactory) *StudentController {
	a.queryRepositoryFactory = f
	return a
}

func (a *StudentController) RegisterRoutes(e *echo.Echo) {
	studentGroup := e.Group("/orders", middleware.Transaction(a.txProvider))
	studentGroup.GET("/get_orders", a.getorders)

}

// @Summary Retrieve all Orders
// @Descript Orders
// @Produce json
// @Tags Orders
// @Router /orders/get_orders [get]
// @Success 200 {object} model.Order
// @Failure 400 {object} model.ErrValidation
func (a *StudentController) getorders(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.Getorders(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}
