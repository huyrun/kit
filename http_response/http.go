package http_response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrInternal     = "có lỗi hệ thống xảy ra"
	ErrUnauthorized = "chưa đăng nhập"
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
	var message string
	switch err.Error() {
	case ErrUnauthorized:
		httpCode = http.StatusUnauthorized
		verdict = VerdictNotAuthorized
		message = ErrUnauthorized
	default:
		httpCode = http.StatusInternalServerError
		verdict = VerdictInternalError
		message = ErrInternal
	}

	Response(c, httpCode, verdict, message, nil)
}

func Abort(c *gin.Context, err error) {
	var httpCode int
	var verdict string
	var message string
	switch err.Error() {
	case ErrUnauthorized:
		httpCode = http.StatusUnauthorized
		verdict = VerdictNotAuthorized
		message = ErrUnauthorized
	default:
		httpCode = http.StatusInternalServerError
		verdict = VerdictInternalError
		message = ErrInternal
	}

	c.AbortWithStatusJSON(httpCode, gin.H{
		"status":  httpCode,
		"verdict": verdict,
		"message": message,
	})
}
