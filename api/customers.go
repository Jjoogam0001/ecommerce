package api

import (
	"context"
	"net/http"
	"strconv"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/repository"
	"dev.azure.com/jjoogam/Ecommerce-core/model"

	"emperror.dev/errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type (
	CustomerController struct {
		queryRepositoryFactory CustomerQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	CustomerQueryRepository interface {
		GetCustomers(ctx context.Context) ([]model.Customer, error)
		FindCustomer(ctx context.Context, customerNumber int) (*model.Customer, error)
		DeleteCustomer(ctx context.Context, customerNumber int) error
	}

	CustomerQueryRepositoryFactory func(pgx.Tx) CustomerQueryRepository
)

func NewCustomerController(txProvider middleware.TxProvider) *CustomerController {
	return &CustomerController{
		queryRepositoryFactory: defaultCustomerQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultCustomerQueryRepositoryFactory(tx pgx.Tx) CustomerQueryRepository {
	return repository.NewCustomerQueryRepository(tx)
}

func (a *CustomerController) WithQueryRepository(f CustomerQueryRepositoryFactory) *CustomerController {
	a.queryRepositoryFactory = f
	return a
}

func (a *CustomerController) RegisterRoutes(e *echo.Echo) {
	customerGroup := e.Group("/customers", middleware.Transaction(a.txProvider))
	customerGroup.GET("/get_customers", a.Getcustomers)
	customerGroup.GET("/get_customer", a.findCustomer)
	customerGroup.DELETE("/delete_customer", a.deleteCustomer)

}

// @Summary Retrieve all customers
// @Description Gets all customers
// @Tags Customer
// @Produce json
// @Router /customers/get_customers [get]
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrValidation
func (a *CustomerController) Getcustomers(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Errorf("unable to resolve transaction", err)
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetCustomers(ctx)
	if err != nil {
		return errors.Errorf("unable to resolve transaction", err)
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve a single customer
// @Description Fetch a single customer
// @Tags Customer
// @Produce json
// @Router /customers/get_customer [get]
// @Param customer_number query int true "customer_number mandatory"
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrValidation
func (a *CustomerController) findCustomer(c echo.Context) error {
	cuid, err := a.decodeCustomer(c)
	if err != nil {
		return errors.Errorf("unable to decode", err)
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Errorf("unable to resolve transaction", err)
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindCustomer(ctx, *cuid)
	if err != nil {
		return errors.Errorf("cant find customer", err)
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary Deletes a single customer
// @Description Deletes a single customer
// @Tags Customer
// @Produce json
// @Router /customers/delete_customer [delete]
// @Param customer_number query int true "customer_number mandatory"
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrValidation
func (a *CustomerController) deleteCustomer(c echo.Context) error {
	cuid, err := a.decodeCustomer(c)
	if err != nil {
		return errors.Errorf("unable to decode", err)
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return errors.Errorf("unable to resolve transaction", err)
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindCustomer(ctx, *cuid)
	if err != nil {
		return errors.Errorf("unable to find customer to delete", err)
	}
	err = r.DeleteCustomer(ctx, *cuid)
	if err != nil {
		return errors.Errorf("cant delete customer", err)
	}

	return c.JSON(http.StatusOK, model.CustomerResponse{
		Customer: *customer,
		Status:   "Deleted",
	})

}

func (h *CustomerController) decodeCustomer(c echo.Context) (*int, error) {
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
