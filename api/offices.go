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
	OfficeController struct {
		queryRepositoryFactory OfficeQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	OfficeQueryRepository interface {
		GetOffices(ctx context.Context) ([]model.Office, error)
		FindOffice(ctx context.Context, officeCode string) (*model.Office, error)
		DeleteOffice(ctx context.Context, officeCode string) error
	}

	OfficeQueryRepositoryFactory func(pgx.Tx) OfficeQueryRepository
)

func NewOfficeController(txProvider middleware.TxProvider) *OfficeController {
	return &OfficeController{
		queryRepositoryFactory: defaultOfficeQueryRepositoryFactory,
		txProvider:             txProvider,
	}
}

func defaultOfficeQueryRepositoryFactory(tx pgx.Tx) OfficeQueryRepository {
	return repository.NewOfficeQueryRepository(tx)
}

func (a *OfficeController) WithQueryRepository(f OfficeQueryRepositoryFactory) *OfficeController {
	a.queryRepositoryFactory = f
	return a
}

func (a *OfficeController) RegisterRoutes(e *echo.Echo) {
	officeGroup := e.Group("/offices", middleware.Transaction(a.txProvider))
	officeGroup.GET("/get_offices", a.GetOffices)
	officeGroup.GET("/get_office", a.findOffice)
	officeGroup.DELETE("/delete_office", a.deleteOffice)

}

// @Summary Retrieve all Offices
// @Description Gets all Offices
// @Tags 	Offices
// @Produce json
// @Router /offices/get_offices [get]
// @Success 200 {object} model.Office
// @Failure 400 {object} model.ErrValidation
func (a *OfficeController) GetOffices(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetOffices(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve a single Office
// @Description Fetch a single Office
// @Tags Offices
// @Produce json
// @Router /offices/get_office [get]
// @Param office_code query string true "office_code mandatory"
// @Success 200 {object} model.Office
// @Failure 400 {object} model.ErrValidation
func (a *OfficeController) findOffice(c echo.Context) error {
	cuid, err := a.decodeOffice(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindOffice(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary Deletes a single Office
// @Description Deletes a single Office
// @Tags Offices
// @Produce json
// @Router /offices/delete_office [delete]
// @Param office_code query string true "office_code mandatory"
// @Success 200 {object} model.Office
// @Failure 400 {object} model.ErrValidation
func (a *OfficeController) deleteOffice(c echo.Context) error {
	cuid, err := a.decodeOffice(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindOffice(ctx, *cuid)

	if err != nil {
		return err
	}

	err = r.DeleteOffice(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)

}

func (h *OfficeController) decodeOffice(c echo.Context) (*string, error) {
	customerNumber := c.QueryParam("office_code")
	if customerNumber == "" && len(customerNumber) == 0 {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "office_code", Reason: "Missing key office_code ."}}}
	}

	return &customerNumber, nil
}
