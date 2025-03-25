package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	apperror "montelukast/pkg/error"
	"montelukast/pkg/wrapper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/go-playground/validator/v10"
)

func ErrorMiddleware(c *gin.Context) {
	c.Next()

	var fieldErrors []apperror.Error

	if len(c.Errors) > 0 {
		var ve validator.ValidationErrors
		if errors.As(c.Errors[0], &ve) {
			for _, fe := range ve {
				fieldErrors = append(fieldErrors, apperror.Error{
					Field:  fe.Field(),
					Detail: apperror.ExtractValidationError(fe),
				})
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, wrapper.ResponseData(nil, "validation error", fieldErrors))
			return
		}

		var je *json.UnmarshalTypeError
		if errors.As(c.Errors[0], &je) {
			jeErrorArray := strings.Split(je.Error(), " ")
			typeError := jeErrorArray[len(jeErrorArray)-1]
			fieldErrors = append(fieldErrors, apperror.Error{
				Field:  strings.ToLower(je.Field),
				Detail: fmt.Sprintf("should be %s data type", typeError),
			})
			c.AbortWithStatusJSON(http.StatusBadRequest, wrapper.ResponseData(nil, "unmarshal type error", fieldErrors))
			return
		}
		var er *apperror.ErrorStruct
		if errors.As(c.Errors[0], &er) {
			fieldError := apperror.Error{
				Field:  er.Field,
				Detail: er.Message,
			}
			c.JSON(er.Status, wrapper.ResponseData(nil, "", []apperror.Error{fieldError}))
			c.Abort()
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "server", "error": "internal server error"})
	}
}
