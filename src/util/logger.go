package util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

const ApiErrorLoggingFormat string = "(%s) -> \"%s\":\nIp: %s\nRoute: %s %s\nStacktrace:\n %s"
const WebsocketErrorLoggingFormat string = "(%s) -> \"%s\":\nIp: %s\nWebsocket type: %s\nStacktrace:\n %s"

func LogApiError(err error, id uuid.UUID, c *gin.Context) {
	log.SetPrefix("[ExceptionHandler] ")

	var stackTrace = ""
	var goError *errors.Error
	if errors.As(err, &goError) {
		stackTrace = goError.ErrorStack()
	}

	log.Printf(
		ApiErrorLoggingFormat,
		id.String(),
		err.Error(),
		c.ClientIP(),
		c.Request.Method,
		c.Request.RequestURI,
		stackTrace,
	)
}

func LogWebsocketError(err error, id uuid.UUID, c *websocket.Conn, t string) {
	log.SetPrefix("[ExceptionHandler] ")

	var stackTrace = ""
	var goError *errors.Error
	if errors.As(err, &goError) {
		stackTrace = goError.ErrorStack()
	}

	log.Printf(
		WebsocketErrorLoggingFormat,
		id.String(),
		err.Error(),
		c.RemoteAddr(),
		t,
		stackTrace,
	)
}
