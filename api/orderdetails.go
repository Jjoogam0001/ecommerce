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
	OrderDetailController struct {
		queryRepositoryFactory OrderDetailQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	OrderDetailQueryRepository interface {
		GetOrderDetails(ctx context.Context) ([]model.OrderDetail, error)
	}

	OrderDetailQueryRepositoryFactory func(pgx.Tx) OrderDetailQueryRepository
)

func NewOrderDetailController(txProvider middleware.TxProvider) *OrderDetailController {
	return &OrderDetailController{
		queryRepositoryFactory: defaultOrderDetailQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultOrderDetailQueryRepositoryFactory(tx pgx.Tx) OrderDetailQueryRepository {
	return repository.NewOrderDetailQueryRepository(tx)
}

func (a *OrderDetailController) WithQueryRepository(f OrderDetailQueryRepositoryFactory) *OrderDetailController {
	a.queryRepositoryFactory = f
	return a
}

func (a *OrderDetailController) RegisterRoutes(e *echo.Echo) {
	orderDetailGroup := e.Group("/orderDetails", middleware.Transaction(a.txProvider))
	orderDetailGroup.GET("/get_orderDetails", a.GetOrderDetails)

}

// @Summary Retrieve all OrderDetails
// @Descript OrderDetails
// @Tags OrderDetails
// @Produce json
// @Router /orderDetails/get_orderDetails [get]
// @Success 200 {object} model.OrderDetail
// @Failure 400 {object} model.ErrValidation
func (a *OrderDetailController) GetOrderDetails(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetOrderDetails(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}

	return c.JSON(http.StatusOK, orders)
}
