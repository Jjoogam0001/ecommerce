package api

import (
	"context"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		UpdatePayment(ctx context.Context, payment model.Payment) error
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
	PaymentGroup.GET("/payments", a.GetPayments)
	PaymentGroup.GET("/get", a.findPayment)
	PaymentGroup.DELETE("/delete", a.deletePayment)
	PaymentGroup.PATCH("/delete", a.UpdatePayment)

}

// @Summary Retrieve all Payments
// @Descript Payments
// @Produce json
// @Tags Payments
// @Router /payments/get [get]
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

// / @Summary Retrieve a  Payment
// @Description Fetch a Payment
// @Tags Payments
// @Produce json
// @Router /payments/get [get]
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
// @Router /payments/delete [delete]
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

// @Summary Updates a  Payment
// @Description updates a single Employee
// @Tags Payments
// @Produce json
// @Router /payments/update [patch]
// @Param payment body model.Payment true "payment_date"
// @Success 200 {object} model.Payment
// @Failure 400 {object} model.ErrValidation
func (a *PaymentController) UpdatePayment(c echo.Context) error {
	ctx := context.Background()
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		metrics.DBErrorInc()
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)
	requestbody, err := ioutil.ReadAll(c.Request().Body)

	var request model.Payment
	err = json.Unmarshal(requestbody, &request)
	if err != nil {
		return err
	}

	err = r.UpdatePayment(ctx, request)
	if err != nil {
		return fmt.Errorf("Error in the request", err.Error())
	}
	return c.String(http.StatusOK, "Customer updated successfully")

}
