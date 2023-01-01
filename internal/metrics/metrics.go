package metrics

import (
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"

	"dev.azure.com/jjoogam/Ecommerce-core/api/middleware"
	"github.com/pkg/errors"

	"github.com/labstack/echo/v4"
)

var host, _ = os.Hostname()

// Prometheus returns prometheus middleware.
func Prometheus() echo.MiddlewareFunc {
	initializeApiMetricsCounters()
	initializeDbMetricsCounters()
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			RequestCount.WithLabelValues(LabelValues(c, err, true)...).Inc()
			RequestDuration.WithLabelValues(LabelValues(c, err, false)...).Observe(float64(Elapsed(start)))

			return err
		}
	}
}

func LabelValues(c echo.Context, err error, addOriginator bool) []string {
	status := c.Response().Status

	if err != nil {
		var httpError *echo.HTTPError
		if errors.As(err, &httpError) {
			status = httpError.Code
		}

		if status == 0 || status == http.StatusOK {
			status = http.StatusInternalServerError
		}
	}
	tags := []string{
		host,                 // host
		strconv.Itoa(status), // code
		c.Request().Method,   // method
		c.Path(),             // endpoint
	}

	if addOriginator {
		tags = append(tags, c.Request().Header.Get(middleware.RequestOriginator)) // originator from header
	}

	return tags
}

// LabelsFromCaller returns hostname from os.Hostname(), and function.
func LabelsFromCaller(skip int) []string {
	pc, _, _, _ := runtime.Caller(skip)
	_, function := path.Split(runtime.FuncForPC(pc).Name())

	return []string{
		host,
		function,
	}
}

// LabelsFromHostName returns hostname from os.Hostname()
func LabelsFromHostName() []string {
	return []string{
		host,
	}
}

// Elapsed time from start as float64 for histogram.
func Elapsed(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}

// DBCallSince records query duration since start time.
func DBCallSince(start time.Time) {
	DbCall.WithLabelValues(LabelsFromCaller(2)...).Observe(float64(Elapsed(start))) //nolint:gomnd
	DbCount.WithLabelValues(LabelsFromCaller(2)...).Inc()
}

// DBErrorInc increases query error counter.
func DBErrorInc() {
	DbError.WithLabelValues(LabelsFromHostName()...).Inc() //nolint:gomnd
}
