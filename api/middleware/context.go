package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

type key string

const (
	trackerKey key = "tracker"
	// HeaderRequestTracker defines the key of request tracker.
	HeaderRequestTracker = "request-tracker"
	// HeaderAuditToken defines the key of audit token.
	HeaderAuditToken = "audit-token"
	// HeaderOriginator defines the key of originator.
	HeaderOriginator = "originator"
	// HeaderRemoteUser defines the key of remote user.
	HeaderRemoteUser = "remote-user"
	// HeaderRoles defines the key of roles.
	HeaderRoles = "roles"
	// HeaderOperatorID defines the key of the operator IDs.
	HeaderOperatorID = "operator-id"
	// HeaderRemoteAddress defines the key of remote address.
	HeaderRemoteAddress = "remote-address"
	// HeaderSessionTracker defines the key of session tracker.
	HeaderSessionTracker = "session-tracker"
	// HeaderDebugTracker defines the key of debug tracker.
	HeaderDebugTracker = "debug-tracker"
	// HeaderOriginatorTracker defines the key of originator tracker.
	HeaderOriginatorTracker = "originator-tracker"
)

// Tracker is the struct of the execution context.
type Tracker struct {
	ExecutionID       string
	AuditToken        string
	Originator        string
	RemoteUser        string
	RemoteAddress     string
	Roles             string
	OperatorIDs       string
	RequestTracker    string //ConversationID
	SessionTracker    string //Another ConversationID
	DebugTracker      string
	OriginatorTracker string
}

func newContext(ctx context.Context, header http.Header) context.Context {
	rid := random.String(32)
	header.Set("execution-id", rid)
	// Used by the official echo.Logger middleware.
	header.Set(echo.HeaderXRequestID, rid)

	t := Tracker{
		ExecutionID:       rid,
		RequestTracker:    header.Get(HeaderRequestTracker),
		SessionTracker:    header.Get(HeaderSessionTracker),
		Originator:        header.Get(HeaderOriginator),
		OriginatorTracker: header.Get(HeaderOriginatorTracker),
		RemoteUser:        header.Get(HeaderRemoteUser),
		RemoteAddress:     header.Get(HeaderRemoteAddress),
		Roles:             header.Get(HeaderRoles),
		OperatorIDs:       header.Get(HeaderOperatorID),
		AuditToken:        header.Get(HeaderAuditToken),
		DebugTracker:      header.Get(HeaderDebugTracker),
	}
	return context.WithValue(ctx, trackerKey, t)
}

// FromContext returns the tracker from the context if exists.
func FromContext(ctx context.Context) (Tracker, bool) {
	t, ok := ctx.Value(trackerKey).(Tracker)
	return t, ok
}

// Context is an echo middleware for parsing the tracking context from headers.
func Context(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r := c.Request()
		h := r.Header
		c.SetRequest(r.WithContext(newContext(r.Context(), h)))
		return next(c)
	}
}
