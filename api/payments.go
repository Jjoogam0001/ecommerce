package api

import (
	"context"
	"net/http"
	"strconv"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/repository"
	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type (
	PaymentController struct {
		queryRepositoryFactory PaymentQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	PaymentQueryRepository interface {
		GetPayments(ctx context.Context) ([]model.Payment, error)
		FindPayment(ctx context.Context, customerNumber int) ([]model.Payment, error)
		DeletePayment(ctx context.Context, customerNumber int) error
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
	PaymentGroup.GET("/get_payment", a.findPayment)
	PaymentGroup.DELETE("/delete_payment", a.deletePayment)

}

// @Summary Retrieve all Payments
// @Descript Payments
// @Produce json
// @Tags Payments
// @Router /payments/get_payments [get]
// @Success 200 {object} model.Payment
// @Failure 400 {object} model.ErrValidation
func (a *PaymentController) GetPayments(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetPayments(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve a  Payment
// @Description Fetch a Payment
// @Tags Payments
// @Produce json
// @Router /payments/get_payment [get]
// @Param customer_number query int true "customer_number mandatory"
// @Success 200 {object} model.Payment
// @Failure 400 {object} model.ErrValidation
func (a *PaymentController) findPayment(c echo.Context) error {
	cuid, err := a.decodePayment(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindPayment(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary deletes a  Payment
// @Description deletes a Payment
// @Tags Payments
// @Produce json
// @Router /payments/delete_payment [delete]
// @Param customer_number query int true "customer_number mandatory"
// @Success 200 {object} model.Payment
// @Failure 400 {object} model.ErrValidation
func (a *PaymentController) deletePayment(c echo.Context) error {
	cuid, err := a.decodePayment(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindPayment(ctx, *cuid)
	if err != nil {
		return err
	}
	err = r.DeletePayment(ctx, *cuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.PaymentResponse{
		Payment: customer[1],
		Status:  "Deleted",
	})

}

func (h *PaymentController) decodePayment(c echo.Context) (*int, error) {
	customerNumber := c.QueryParam("customer_number")
	if customerNumber == "" && len(customerNumber) == 0 {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "customer_number", Reason: "Missing key customer_number ."}}}
	}

	cuId, err := strconv.Atoi(customerNumber)
	if err != nil {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "customer_number", Reason: "Incorrect customer_number"}}}
	}

	return &cuId, nil
}
