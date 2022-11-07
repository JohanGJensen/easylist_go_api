package validation

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Errors []ErrorMessage
}

type ErrorMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fe.Param() + " is below minimum allowed characters"
	case "max":
		return fe.Param() + " is above maximum allowed characters"
	}
	return "Unknown error"
}

func Validate(c *gin.Context, e error) {
	ve := validator.ValidationErrors{}

	if errors.As(e, &ve) {
		out := make([]ErrorMessage, len(ve))

		for i, fe := range ve {
			out[i] = ErrorMessage{fe.Field(), getErrorMsg(fe)}
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Errors: out,
		})
	}
}
