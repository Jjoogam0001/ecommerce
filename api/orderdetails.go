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
	OrderDetailController struct {
		queryRepositoryFactory OrderDetailQueryRepositoryFactory
		txProvider             middleware.TxProvider
	}

	OrderDetailQueryRepository interface {
		GetOrderDetails(ctx context.Context) ([]model.OrderDetail, error)
		FindOrderDetails(ctx context.Context, orderNumber int) ([]model.OrderDetail, error)
		DeleteOrder(ctx context.Context, orderNumber int) error
		UpdateOrderDetails(ctx context.Context, ordersDetails model.OrderDetail) error
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
	orderDetailGroup.GET("/get_orderDetail", a.findOrderDetails)
	orderDetailGroup.DELETE("/delete_orderDetail", a.deleteOrderDetails)
	orderDetailGroup.PATCH("/update", a.UpdateOrderDetails)

}

// @Summary Retrieve all OrderDetails
// @Descript OrderDetails
// @Tags OrderDetails
// @Produce json
// @Router /orderDetails/get [get]
// @Success 200 {object} model.OrderDetail
// @Failure 400 {object} model.ErrValidation
func (a *OrderDetailController) GetOrderDetails(c echo.Context) error {

	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)

	ctx := c.Request().Context()
	orders, err := r.GetOrderDetails(ctx)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, orders)
}

// @Summary Retrieve order details
// @Description Fetch order details
// @Tags OrderDetails
// @Produce json
// @Router /orderDetails/get_orderDetail [get]
// @Param order_number query int true "order_number mandatory"
// @Success 200 {object} model.OrderDetail
// @Failure 400 {object} model.ErrValidation
func (a *OrderDetailController) findOrderDetails(c echo.Context) error {
	cuid, err := a.decodeOrderDetails(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindOrderDetails(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, customer)

}

// @Summary Delete order details
// @Description Delete order details
// @Tags Orders
// @Produce json
// @Router /orderDetails/delete [delete]
// @Param order_number query int true "order_number mandatory"
// @Success 200 {object} model.Order
// @Failure 400 {object} model.ErrValidation
func (a *OrderDetailController) deleteOrderDetails(c echo.Context) error {
	cuid, err := a.decodeOrderDetails(c)
	if err != nil {
		return err
	}
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		return err
	}
	r := a.queryRepositoryFactory(db)
	ctx := c.Request().Context()
	customer, err := r.FindOrderDetails(ctx, *cuid)
	if err != nil {
		return err
	}

	err = r.DeleteOrder(ctx, *cuid)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, model.OrderDetailResponse{
		OrderDetail: customer[1],
		Status:      "Deleted",
	})

}

func (h *OrderDetailController) decodeOrderDetails(c echo.Context) (*int, error) {
	customerNumber := c.QueryParam("order_number")
	if customerNumber == "" && len(customerNumber) == 0 {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "order_number", Reason: "Missing key customer_number ."}}}
	}

	cuId, err := strconv.Atoi(customerNumber)
	if err != nil {
		return nil, model.ErrValidation{InvalidParams: []model.InvalidParam{{Name: "order_number", Reason: "Incorrect customer_number"}}}
	}

	return &cuId, nil
}

// @Summary Updates a  OrderDetail
// @Description updates OrderDetail
// @Tags OrderDetails
// @Produce json
// @Router /orderDetails/update [patch]
// @Param customer body model.OrderDetail true "order_number"
// @Success 200 {object} model.OrderDetail
// @Failure 400 {object} model.ErrValidation
func (a *OrderDetailController) UpdateOrderDetails(c echo.Context) error {
	ctx := context.Background()
	db, err := middleware.FromTransactionContext(c)
	if err != nil {
		metrics.DBErrorInc()
		return errors.Wrap(err, "unable to resolve transaction")
	}
	r := a.queryRepositoryFactory(db)
	requestbody, err := ioutil.ReadAll(c.Request().Body)

	var request model.OrderDetail
	err = json.Unmarshal(requestbody, &request)
	if err != nil {
		return err
	}

	err = r.UpdateOrderDetails(ctx, request)
	if err != nil {
		return fmt.Errorf("Error in the request", err.Error())
	}
	return c.String(http.StatusOK, "OrderDetail updated successfully")

}
