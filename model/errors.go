package model

import (
	"encoding/json"
	"net/http"
)

type (
	// ErrValidation defines an error that will be returned in the case of an invalid combination of event-market-outcome IDs.
	ErrValidation struct {
		InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
	}

	// InvalidParam defines the structure of an invalid parameter passed to the risk API.
	InvalidParam struct {
		Name   string `json:"name"`
		Reason string `json:"reason"`
	}

	// Invalid defines behaviour when inputs invalid.
	Invalid interface {
		error
		Invalid() bool
	}
)

// Error interface implementation
func (err ErrValidation) Error() string {
	var errMessage string

	if len(err.InvalidParams) > 0 {
		bytes, _ := json.Marshal(err)
		errMessage = string(bytes)
	}

	return errMessage
}

// Invalid interface implemented
func (err ErrValidation) Invalid() bool {
	return true
}

// IsErrInvalid returns true if the error is caused with an
// object (image, container, network, volume, …) is not found in the docker host.
func IsErrInvalid(err error) bool {
	te, ok := err.(Invalid)
	return ok && te.Invalid()
}

type (
	// ErrConflict defines an error that will be returned if the resource already exists.
	ErrConflict struct {
		Message string `json:"message"`
	}

	// Conflict defines behaviour when resource already exists.
	Conflict interface {
		error
		Conflict() bool
	}
)

// Error interface implementation
func (err ErrConflict) Error() string {
	return err.Message
}

// Conflict interface implementation
func (err ErrConflict) Conflict() bool {
	return true
}

// IsErrConflict returns true if the error is caused by an ErrConflict.
func IsErrConflict(err error) bool {
	te, ok := err.(Conflict)
	return ok && te.Conflict()
}

// NewErrConflict returns a new not found error.
func NewErrConflict(message string) ErrConflict {
	if message == "" {
		message = http.StatusText(http.StatusConflict)
	}

	return ErrConflict{message}
}
