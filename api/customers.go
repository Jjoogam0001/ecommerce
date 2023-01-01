package api

import (
	"context"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/repository"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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
		FindCustomer(ctx context.Context, email string) (*model.Customer, error)
		DeleteCustomer(ctx context.Context, email string) error
		UpdateCustomer(ctx context.Context, customer model.Customer) error
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
	customerGroup.GET("/customers", a.Getcustomers)
	customerGroup.GET("/get", a.findCustomer)
	customerGroup.DELETE("/delete", a.deleteCustomer)
	customerGroup.PATCH("/update", a.UpdateCustomer)

}

// @Summary Retrieve all customers
// @Description Gets all customers
// @Tags Customer
// @Produce json
// @Router /customers/customers [get]
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrValidation
func (a *CustomerController) Getcustomers(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetCustomers(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve a single customer
// @Description Fetch a single customer
// @Tags Customer
// @Produce json
// @Router /customers/get [get]
// @Param email query string true "email mandatory"
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrValidation
func (a *CustomerController) findCustomer(c echo.Context) error {
	email, err := a.decodeCustomer(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindCustomer(ctx, email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary Deletes a single customer
// @Description Deletes a single customer
// @Tags Customer
// @Produce json
// @Router /customers/delete [delete]
// @Param email query string true "email mandatory"
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrValidation
func (a *CustomerController) deleteCustomer(c echo.Context) error {
	email, err := a.decodeCustomer(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	_, err = r.FindCustomer(ctx, email)
	if err != nil {
		return err
	}
	err = r.DeleteCustomer(ctx, email)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.CustomerResponse{
		Status: "Deleted",
	})

}

func (h *CustomerController) decodeCustomer(c echo.Context) (string, error) {
	customerNumber := c.QueryParam("email")
	if customerNumber == "" && len(customerNumber) == 0 {
		return "", model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "customer_number", Reason: "Missing key customer_number ."}}}
	}

	return customerNumber, nil
}

// @Summary Updates a  customer
// @Description updates a single customer
// @Tags Customer
// @Produce json
// @Router /customers/update [patch]
// @Param customer body model.Customer true "customer_number"
// @Success 200 {object} model.Customer
// @Failure 400 {object} model.ErrValidation
func (a *CustomerController) UpdateCustomer(c echo.Context) error {
	ctx := context.Background()
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		metrics.DBErrorInc()
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)
	requestbody, err := ioutil.ReadAll(c.Request().Body)

	var request model.Customer
	err = json.Unmarshal(requestbody, &request)
	if err != nil {
		return err
	}
	if err != nil {

		return errors.Wrap(err, "error getting customer information from the request")
	}

	if request.Email == "" {
		return model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "email", Reason: "email is a mandatory attribute and cannot be 0"}}}
	}
	err = r.UpdateCustomer(ctx, request)
	if err != nil {
		return fmt.Errorf("Error in the request", err.Error())
	}
	return c.String(http.StatusOK, "Customer updated successfully")

}
