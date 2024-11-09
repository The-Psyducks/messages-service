package modelErrors

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SendErrorMessage(ctx *gin.Context, err *MessageError) {
	//ctx.JSON(500, gin.H{"TODO: implement nice errors": err.Error()})
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(err.Status, err)
}

// rfc nomeacuerdocuanto
type MessageError struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

func (m MessageError) Error() string {
	return m.Detail
}

func ValidationError(detail string) *MessageError {
	return &MessageError{
		Type:   "about:blank",
		Title:  "Validation Error",
		Status: http.StatusBadRequest,
		Detail: detail,
	}
}

func ExternalServiceError(detail string) *MessageError {
	return &MessageError{
		Type:   "about:blank",
		Title:  "External Service Error",
		Status: http.StatusServiceUnavailable,
		Detail: detail,
	}
}

func InternalServerError(detail string) *MessageError {
	return &MessageError{
		Type:   "about:blank",
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: detail,
	}
}

func AuthenticationError(details string) *MessageError {
	return &MessageError{
		Type:   "about:blank",
		Title:  "Authentication Error",
		Status: http.StatusUnauthorized,
		Detail: details,
	}
}

func BadRequestError(details string) *MessageError {
	return &MessageError{
		Type:   "about:blank",
		Title:  "Bad Request",
		Status: http.StatusBadRequest,
		Detail: details,
	}
}
