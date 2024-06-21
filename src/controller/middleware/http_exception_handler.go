package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/util"
	"net/http"
)

const UnexpectedErrorTitle string = "Internal server error"
const UnexpectedErrorDescription string = "An unexpected error has occurred"
const DeviceUnreachableTitle string = "Device unreachable"
const DeviceUnreachableDescription string = "The device selected was not able to receive the command"
const InvalidJsonTitle string = "Invalid json"
const InvalidJsonDescription string = "The json provided is invalid"
const ObjectNotFoundTitle string = "Object not found"
const ObjectNotFoundDescription string = "The object requested was not found on the server"

func ExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		var err = c.Errors[0]
		id := uuid.New()
		util.LogApiError(err, id, c)
		var deviceUnreachableError *exceptions.DeviceUnreachable
		if errors.As(err, &deviceUnreachableError) {
			handleDeviceUnreachable(c, id, err.Error())
			return
		}
		var jsonTypeError *json.UnmarshalTypeError
		var jsonSyntaxError *json.SyntaxError
		if errors.As(err, &jsonTypeError) || errors.As(err, &jsonSyntaxError) {
			handleInvalidJson(c, id, err.Error())
			return
		}
		var objectNotFoundError *exceptions.ObjectNotFound
		if errors.As(err, &objectNotFoundError) {
			handleObjectNotFound(c, id, err.Error())
			return
		}
		handleUnexpectedError(c, id)
	}
}

func handleUnexpectedError(c *gin.Context, id uuid.UUID) {
	var err api.ErrorResponse

	err.SetId(id.String())
	err.SetTitle(UnexpectedErrorTitle)
	err.SetStatus(http.StatusInternalServerError)
	err.SetDescription(UnexpectedErrorDescription)
	err.SetExpected(false)
	c.AbortWithStatusJSON(http.StatusInternalServerError, err)
}

func handleDeviceUnreachable(c *gin.Context, id uuid.UUID, message string) {
	var err api.ErrorResponse

	err.SetId(id.String())
	err.SetTitle(DeviceUnreachableTitle)
	err.SetStatus(http.StatusServiceUnavailable)
	err.SetDescription(DeviceUnreachableDescription)
	err.SetMessage(message)
	err.SetExpected(true)
	c.AbortWithStatusJSON(http.StatusServiceUnavailable, err)
}

func handleInvalidJson(c *gin.Context, id uuid.UUID, message string) {
	var err api.ErrorResponse

	err.SetId(id.String())
	err.SetTitle(InvalidJsonTitle)
	err.SetStatus(http.StatusBadRequest)
	err.SetDescription(InvalidJsonDescription)
	err.SetMessage(message)
	err.SetExpected(true)
	c.AbortWithStatusJSON(http.StatusBadRequest, err)
}

func handleObjectNotFound(c *gin.Context, id uuid.UUID, message string) {
	var err api.ErrorResponse

	err.SetId(id.String())
	err.SetTitle(ObjectNotFoundTitle)
	err.SetStatus(http.StatusNotFound)
	err.SetDescription(ObjectNotFoundDescription)
	err.SetMessage(message)
	err.SetExpected(true)
	c.AbortWithStatusJSON(http.StatusNotFound, err)
}
