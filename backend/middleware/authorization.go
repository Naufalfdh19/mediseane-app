package middleware

import (
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"

	"github.com/gin-gonic/gin"
)

func AuthUserMiddleware(c *gin.Context) {
	role, found := c.Get("role")
	if !found {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrCredentialWrong, nil))
		c.Abort()
		return
	}
	if role != appconstant.ROLE_USER {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrCredentialWrong, nil))
		c.Abort()
		return
	}
	c.Next()
}

func AuthPharmacistMiddleware(c *gin.Context) {
	role, found := c.Get("role")
	if !found {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrCredentialWrong, nil))
		c.Abort()
		return
	}
	if role != appconstant.ROLE_PHARMACY {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrCredentialWrong, nil))
		c.Abort()
		return
	}
	c.Next()
}

func AuthAdminMiddleware(c *gin.Context) {
	role, found := c.Get("role")
	if !found {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrCredentialWrong, nil))
		c.Abort()
		return
	}
	if role != appconstant.ROLE_ADMIN {
		c.Error(apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrCredentialWrong, nil))
		c.Abort()
		return
	}
	c.Next()
}
