package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// RequestOriginator gets the name used in the headers and logs.
	RequestOriginator = "Originator"
)

// Originator validates that the request contains a Request Originator field in the request header.
func Originator() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Path(), "/metrics") ||
				strings.Contains(c.Path(), "/swagger") ||
				strings.Contains(c.Path(), "/healthcheck") {
				return next(c)
			}
			rt := c.Request().Header.Get(RequestOriginator)
			if rt == "" {
				c.Request().Header.Set(RequestOriginator, "N/A")
			}

			return next(c)
		}
	}
}
