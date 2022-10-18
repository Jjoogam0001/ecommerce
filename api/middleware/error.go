package middleware

import (
	"net/http"

	"dev.azure.com/jjoogam/Ecommerce-core/model"
	"github.com/labstack/echo/v4"
)

// Error middleware to handle errors and their http status codes.
func Error(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			switch {
			case model.IsErrInvalid(err):
				return echo.NewHTTPError(http.StatusBadRequest, err)
			case model.IsErrConflict(err):
				return echo.NewHTTPError(http.StatusConflict, err)
			default:
				return echo.NewHTTPError(http.StatusInternalServerError, err)
			}
		}

		return nil
	}
}
