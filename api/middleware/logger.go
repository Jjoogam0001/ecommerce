package middleware

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/gommon/log"
)

var (
	defaultSkipper = func(c echo.Context) bool {
		return strings.HasPrefix(c.Request().RequestURI, "/healthcheck")
	}

	defaultLoggerConfig = middleware.LoggerConfig{
		Skipper: defaultSkipper,
		Format: fmt.Sprintf(`{"time":"${time_rfc3339_nano}",`+
			`"execution_id":"${id}","request_tracker":"${header:%v}","session_tracker":"${header:%v}",`+
			`"originator":"${header:%v}","remote_address":"${header:%v}","operator_id":"${header:%v}",`+
			`"remote_ip":"${remote_ip}","host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",`+
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"`+
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}`+"\n",
			HeaderRequestTracker, HeaderSessionTracker, HeaderOriginator, HeaderRemoteAddress, HeaderOperatorID),
		CustomTimeFormat: middleware.DefaultLoggerConfig.CustomTimeFormat,
	}
)

// UseLogger configures the log level, the context tracker and the logging middleware.
// It build a logger instance with the tracking context and this instance is attached into the current request.
// In that way every line of code written with this log instance, will print the execution and correlation id.
// Then we will have a more standard logging and helps to troubleshooting.
func UseLogger(e *echo.Echo) {
	if log.Level() == log.DEBUG {
		useBodyDump(e)
	}

	e.Logger.SetLevel(log.Level())
	e.Use(Context)
	e.Use(middleware.LoggerWithConfig(defaultLoggerConfig))
	e.Use(Logger)
}

// Logger is an echo middleware which configure the logger instances with the execution id and the correlation id.
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		t, ok := FromContext(c.Request().Context())
		if !ok {
			return next(c)
		}
		echoLogger := c.Logger()
		prefix := fmt.Sprintf(`%v","execution_id":"%v","request_tracker":"%v","session_tracker":"%v`, echoLogger.Prefix(), t.ExecutionID, t.RequestTracker, t.SessionTracker)
		l := log.New(prefix)
		l.SetLevel(echoLogger.Level())
		l.SetOutput(echoLogger.Output())
		c.SetLogger(l)
		return next(c)
	}
}

func useBodyDump(e *echo.Echo) {
	e.Use(middleware.BodyDumpWithConfig(
		middleware.BodyDumpConfig{
			Skipper: defaultSkipper,
			Handler: func(c echo.Context, req, resp []byte) {
				sreq, _ := url.QueryUnescape(string(req[:]))
				sresp, _ := url.QueryUnescape(string(resp[:]))

				c.Logger().Debugj(log.JSON{
					"uri":      c.Request().RequestURI,
					"method":   c.Request().Method,
					"request":  sreq,
					"response": sresp,
				})
			},
		},
	))
}
