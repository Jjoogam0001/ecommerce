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
	PaymentController struct {
		queryRepositoryFactory PaymentQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	PaymentQueryRepository interface {
		GetPayments(ctx context.Context) ([]model.Payment, error)
	}

	PaymentQueryRepositoryFactory func(pgx.Tx) PaymentQueryRepository
)

func NewPaymentController(txProvider middleware.TxProvider) *PaymentController {
	return &PaymentController{
		queryRepositoryFactory: defaultPaymentQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultPaymentQueryRepositoryFactory(tx pgx.Tx) PaymentQueryRepository {
	return repository.NewPaymentQueryRepository(tx)
}

func (a *PaymentController) WithQueryRepository(f PaymentQueryRepositoryFactory) *PaymentController {
	a.queryRepositoryFactory = f
	return a
}

func (a *PaymentController) RegisterRoutes(e *echo.Echo) {
	PaymentGroup := e.Group("/payments", middleware.Transaction(a.txProvider))
	PaymentGroup.GET("/get_payments", a.GetPayments)

}

// @Summary Retrieve all Payments
// @Descript Payments
// @Produce json
// @Tags Payment
// @Router /payments/get_payments [get]
// @Success 200 {object} model.Payment
// @Failure 400 {object} model.ErrValidation
func (a *PaymentController) GetPayments(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetPayments(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}

	return c.JSON(http.StatusOK, orders)
}
