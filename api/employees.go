package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/metrics"
	"dev.azure.com/jjoogam/Ecommerce-core/internal/repository"
	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"emperror.dev/errors"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
)

type (
	EmployeeController struct {
		queryRepositoryFactory EmployeeQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	EmployeeQueryRepository interface {
		GetEmployees(ctx context.Context) ([]model.Employee, error)
		FindEmployee(ctx context.Context, employeeNumber int) (*model.Employee, error)
		DeleteEmployee(ctx context.Context, employeeNumber int) error
		UpdateEmployee(ctx context.Context, employee model.Employee) error
	}

	EmployeeQueryRepositoryFactory func(pgx.Tx) EmployeeQueryRepository
)

func NewEmployeeController(txProvider middleware.TxProvider) *EmployeeController {
	return &EmployeeController{
		queryRepositoryFactory: defaultEmployeeQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultEmployeeQueryRepositoryFactory(tx pgx.Tx) EmployeeQueryRepository {
	return repository.NewEmployeeQueryRepository(tx)
}

func (a *EmployeeController) WithQueryRepository(f EmployeeQueryRepositoryFactory) *EmployeeController {
	a.queryRepositoryFactory = f
	return a
}

func (a *EmployeeController) RegisterRoutes(e *echo.Echo) {
	employeeGroup := e.Group("/employees", middleware.Transaction(a.txProvider))
	employeeGroup.GET("/employees", a.GetEmployees)
	employeeGroup.GET("/get", a.findemployee)
	employeeGroup.DELETE("/delete", a.deleteEmployee)
	employeeGroup.PATCH("/update", a.UpdateEmployee)

}

// @Summary Retrieve all Employees
// @Description Gets all Employees
// @Tags Employees
// @Produce json
// @Router /employees/employees [get]
// @Success 200 {object} model.Employee
// @Failure 400 {object} model.ErrValidation
func (a *EmployeeController) GetEmployees(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetEmployees(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve a single Employee
// @Description Fetch a single Employee
// @Tags Employees
// @Produce json
// @Router /employees/get [get]
// @Param employee_number query int true "employee_number mandatory"
// @Success 200 {object} model.Employee
// @Failure 400 {object} model.ErrValidation
func (a *EmployeeController) findemployee(c echo.Context) error {
	cuid, err := a.decodeEmployee(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindEmployee(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary deletes a single Employee
// @Description deleted a single Employee
// @Tags Employees
// @Produce json
// @Router /employees/delete [delete]
// @Param employee_number query int true "employee_number mandatory"
// @Success 200 {object} model.Employee
// @Failure 400 {object} model.ErrValidation
func (a *EmployeeController) deleteEmployee(c echo.Context) error {
	cuid, err := a.decodeEmployee(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindEmployee(ctx, *cuid)
	if err != nil {
		return err
	}
	err = r.DeleteEmployee(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.EmployeeResponse{
		Employee: *customer,
		Status:   "Delete",
	})

}

func (h *EmployeeController) decodeEmployee(c echo.Context) (*int, error) {
	customerNumber := c.QueryParam("employee_number")
	if customerNumber == "" && len(customerNumber) == 0 {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "employee_number", Reason: "Missing key customer_number ."}}}
	}

	cuId, err := strconv.Atoi(customerNumber)
	if err != nil {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "employee_number", Reason: "Incorrect customer_number"}}}
	}

	return &cuId, nil
}

// @Summary Updates a  Employee
// @Description updates a single Employee
// @Tags Employees
// @Produce json
// @Router /employees/update [patch]
// @Param customer body model.Employee true "employee_number"
// @Success 200 {object} model.Employee
// @Failure 400 {object} model.ErrValidation
func (a *EmployeeController) UpdateEmployee(c echo.Context) error {
	ctx := context.Background()
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		metrics.DBErrorInc()
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)
	requestbody, err := ioutil.ReadAll(c.Request().Body)

	var request model.Employee
	err = json.Unmarshal(requestbody, &request)
	if err != nil {
		return err
	}
	if err != nil {

		return errors.Wrap(err, "error getting employee information from the request")
	}

	if request.Email == "" {
		return model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "email", Reason: "email is a mandatory attribute and cannot be 0"}}}
	}
	err = r.UpdateEmployee(ctx, request)
	if err != nil {
		return fmt.Errorf("Error in the request", err.Error())
	}
	return c.String(http.StatusOK, "Employee updated successfully")

}
