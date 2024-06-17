package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"log"
)

const ErrorLoggingFormat string = "(%s) -> \"%s\":\nIp: %s\nRoute: %s %s\nStacktrace:\n %s"

func LogError(err error, id uuid.UUID, c *gin.Context) {
	log.SetPrefix("[ExceptionHandler] ")

	var stackTrace = ""
	var goError *errors.Error
	if errors.As(err, &goError) {
		stackTrace = goError.ErrorStack()
	}

	log.Printf(
		ErrorLoggingFormat,
		id.String(),
		err.Error(),
		c.ClientIP(),
		c.Request.Method,
		c.Request.RequestURI,
		stackTrace,
	)
}
