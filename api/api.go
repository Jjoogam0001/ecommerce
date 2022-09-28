package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	echoSwagger "github.com/swaggo/echo-swagger"

	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/api/docs/swagger"
	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/api/middleware"
	"dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/config"
	_ "dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/docs/swagger"
	_ "dev.azure.com/jjoogam0290/HelloWorld/HelloWorld/model"
)

// API defines the functions of the risk HTTP server.
type (
	API struct {
		server *echo.Echo
	}

	Controller interface {
		RegisterRoutes(e *echo.Echo)
	}
) // Martsoft Inc E-Commerce API.

// @title Martsoft Inc E-Commerce API 2.0
// @version 1.0
// @description This is the API E-Commerce businesses.

// @contact.name Martsoft Inc
// @contact.email Martsoftfilmz@gmail.com

// @schemes http
func NewAPI(c config.AppConfig) *API {
	e := echo.New()

	middleware.UseLogger(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/healthcheck", health)
	initSwagerInfo(c)
	return &API{server: e}
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK\n")
}

// UsingControllers registers in the server the controllers passed by params.
func (a *API) UsingControllers(controllers ...Controller) *echo.Echo {
	for _, controller := range controllers {
		controller.RegisterRoutes(a.server)
	}
	return a.server
}

// UsingDefaultControllers registers in the server the default controllers configuration needed for running the server.
func (a *API) UsingDefaultControllers(txProvider middleware.TxProvider) *echo.Echo {
	return a.UsingControllers([]Controller{
		NewStudentController(txProvider),
		NewCustomerController(txProvider),
		NewEmployeeController(txProvider),
		NewOfficeController(txProvider),
		NewOrderDetailController(txProvider),
		NewPaymentController(txProvider),
		NewProductController(txProvider),
		NewProductLineController(txProvider),
	}...)
}

// Close closes a risk HTTP server.
func (a *API) Close() error {
	return a.server.Close()
}
func initSwagerInfo(c config.AppConfig) {
	swagger.SwaggerInfo.Title = "Martsoft Inc E-Commerce API 2.0"
	swagger.SwaggerInfo.Description = "This is the API E-Commerce businesses."
	swagger.SwaggerInfo.Version = "1.0"
	swagger.SwaggerInfo.Host = c.Swagger.Host
	swagger.SwaggerInfo.BasePath = ""
	swagger.SwaggerInfo.Schemes = []string{"http", "https"}
}
