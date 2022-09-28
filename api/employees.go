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
	EmployeeController struct {
		queryRepositoryFactory EmployeeQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	EmployeeQueryRepository interface {
		GetEmployees(ctx context.Context) ([]model.Employee, error)
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
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetEmployees(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to resolve transaction")
	}

	return c.JSON(http.StatusOK, orders)
}
