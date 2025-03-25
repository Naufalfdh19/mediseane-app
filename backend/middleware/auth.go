package middleware

import (
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"
	jwttoken "montelukast/pkg/jwt_token"
	"strings"

	"github.com/gin-gonic/gin"
)

func CheckAuthorization(c *gin.Context) {
	var tokenString string
	authHeader := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeader) <= 1 {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrTokenInvalid, apperror.ErrTokenInvalid)
		c.Error(err)
		c.Abort()
		return
	}
	tokenString = authHeader[1]

	jwt_token := jwttoken.JwtTokenImpl{}
	jwtTokenClaims, err := jwt_token.ParseJwtTokenForAuth(c, tokenString)
	if err != nil {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, err)
		c.Error(err)
		c.Abort()
		return
	}
	
	if jwtTokenClaims.UserID == "0" || jwtTokenClaims.Type != appconstant.JwtTokenAuthType {
		err := apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUserUnauthorized, nil)
		c.Error(err)
		c.Abort()
		return
	}

	c.Set("user_id", jwtTokenClaims.UserID)
	c.Set("type", jwtTokenClaims.Type)
	c.Set("role", jwtTokenClaims.Role)
	c.Next()
}

func CheckProductAuthorization(c *gin.Context) {
	var tokenString string
	authHeader := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeader) <= 1 {
		c.Set("user_id", "-1")
		return
	}
	tokenString = authHeader[1]

	jwt_token := jwttoken.JwtTokenImpl{}
	jwtTokenClaims, err := jwt_token.ParseJwtTokenForAuth(c, tokenString)
	if err != nil {
		c.Set("user_id", "-1")
		return
	}
	
	if jwtTokenClaims.UserID == "0" || jwtTokenClaims.Type != appconstant.JwtTokenAuthType {
		c.Set("user_id", "-1")
		return
	}

	c.Set("user_id", jwtTokenClaims.UserID)
	c.Next()
}
