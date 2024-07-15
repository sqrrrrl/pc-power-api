package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/exceptions"
	"github.com/pc-power-api/src/util"
	"net/http"
	"strings"
)

const UnexpectedErrorTitle string = "Internal server error"
const UnexpectedErrorDescription string = "An unexpected error has occurred"
const DeviceUnreachableTitle string = "Device unreachable"
const DeviceUnreachableDescription string = "The device selected was not able to receive the command"
const InvalidJsonTitle string = "Invalid json"
const InvalidJsonDescription string = "The json provided is invalid"
const ObjectNotFoundTitle string = "Object not found"
const ObjectNotFoundDescription string = "The object requested was not found on the server"
const NoAccessTitle string = "No access"
const NoAccessDescription string = "The user does not have access to this resource"
const ValidationErrorTitle string = "Validation error"
const ValidationErrorDescription string = "The input provided is invalid"

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
		var noAccessError *exceptions.NoAccess
		if errors.As(err, &noAccessError) {
			handleNoAccess(c, id, err.Error())
			return
		}
		var validationError validator.ValidationErrors
		if errors.As(err, &validationError) {
			handleValidationErrors(c, id, validationError)
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

func handleNoAccess(c *gin.Context, id uuid.UUID, message string) {
	var err api.ErrorResponse

	err.SetId(id.String())
	err.SetTitle(NoAccessTitle)
	err.SetStatus(http.StatusForbidden)
	err.SetDescription(NoAccessDescription)
	err.SetMessage(message)
	err.SetExpected(true)
	c.AbortWithStatusJSON(http.StatusForbidden, err)
}

func handleValidationErrors(c *gin.Context, id uuid.UUID, validationErrors validator.ValidationErrors) {
	var err api.ErrorResponse

	err.SetId(id.String())
	err.SetTitle(ValidationErrorTitle)
	err.SetStatus(http.StatusBadRequest)
	err.SetDescription(ValidationErrorDescription)
	err.SetErrors(translateValidationErrors(validationErrors))
	err.SetExpected(true)
	c.AbortWithStatusJSON(http.StatusBadRequest, err)
}

func translateValidationErrors(validationErrors validator.ValidationErrors) []string {
	var translatedErrors []string
	for _, validationError := range validationErrors {
		translatedError := validationError.Error()
		switch validationError.Tag() {
		case "required":
			translatedError = validationError.Field() + " is required"
		case "len":
			translatedError = validationError.Field() + " must be " + validationError.Param() + " characters long"
		case "uuid":
			translatedError = validationError.Field() + " must be a valid uuid"
		case "max":
			translatedError = validationError.Field() + " must be at most " + validationError.Param() + " characters long"
		case "min":
			translatedError = validationError.Field() + " must be at least " + validationError.Param() + " characters long"
		case "eqfield":
			translatedError = validationError.Field() + " must be equal to " + validationError.Param()
		case "excludesall":
			if validationError.Param() == " " {
				translatedError = validationError.Field() + " must not contain spaces"
			} else {
				translatedError = validationError.Field() + " must not contain " + validationError.Param()
			}
		case "printascii":
			translatedError = validationError.Field() + " must only contain printable ascii characters"
		case "oneof":
			translatedError = validationError.Field() + " must be one of " + strings.Replace(validationError.Param(), " ", ", ", -1)
		}
		translatedErrors = append(translatedErrors, translatedError)
	}
	return translatedErrors
}
