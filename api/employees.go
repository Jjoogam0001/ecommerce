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
	EmployeeController struct {
		queryRepositoryFactory EmployeeQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	EmployeeQueryRepository interface {
		GetEmployees(ctx context.Context) ([]model.Employee, error)
		FindEmployee(ctx context.Context, employeeNumber int) (*model.Employee, error)
		DeleteEmployee(ctx context.Context, employeeNumber int) error
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
	employeeGroup.GET("/get_employees", a.GetEmployees)
	employeeGroup.GET("/get_employee", a.findemployee)
	employeeGroup.DELETE("/delete_employee", a.deleteEmployee)

}

// @Summary Retrieve all Employees
// @Description Gets all Employees
// @Tags Employees
// @Produce json
// @Router /employees/get_employees [get]
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
// @Router /employees/get_employee [get]
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
// @Router /employees/delete_employee [delete]
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
