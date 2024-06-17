package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pc-power-api/src/api"
	"github.com/pc-power-api/src/util"
	"net/http"
)

const UnexpectedErrorTitle string = "Internal server error"
const UnexpectedErrorDescription string = "An unexpected error has occurred"

func ExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			id := uuid.New()
			util.LogError(err, id, c)
			handleUnexpectedError(c, id)
		}
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
