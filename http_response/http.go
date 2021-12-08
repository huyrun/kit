package http_response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, status int, verdict string, message string, data interface{}) {
	c.JSON(status, gin.H{
		"verdict": verdict,
		"message": message,
		"data":    data,
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	Response(c, http.StatusOK, VerdictSuccess, message, data)
}

func BadRequest(c *gin.Context, message string, data interface{}) {
	Response(c, http.StatusBadRequest, VerdictBadRequest, message, data)
}

func NotAuthorized(c *gin.Context, message string) {
	Response(c, http.StatusUnauthorized, VerdictNotAuthorized, message, nil)
}

func Error(c *gin.Context, err error) {
	var httpCode int
	var verdict string

	var errUnauthorized *ErrUnauthorized
	switch {
	case errors.As(err, &errUnauthorized):
		httpCode = http.StatusUnauthorized
		verdict = VerdictNotAuthorized
	default:
		httpCode = http.StatusInternalServerError
		verdict = VerdictInternalError
	}

	Response(c, httpCode, verdict, err.Error(), nil)
}

func Abort(c *gin.Context, err error) {
	var httpCode int
	var verdict string

	var errUnauthorized *ErrUnauthorized
	switch {
	case errors.As(err, &errUnauthorized):
		httpCode = http.StatusUnauthorized
		verdict = VerdictNotAuthorized
	default:
		httpCode = http.StatusInternalServerError
		verdict = VerdictInternalError
	}

	c.AbortWithStatusJSON(httpCode, gin.H{
		"status":  httpCode,
		"verdict": verdict,
		"message": err.Error(),
	})
}
