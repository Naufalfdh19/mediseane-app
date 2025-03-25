package jwttoken

import (
	"context"
	appconstant "montelukast/pkg/constant"
	apperror "montelukast/pkg/error"

	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenItf interface {
	GenerateJwtToken(role, userID string) (string, error)
	ParseJwtToken(*gin.Context, string) (*jwt.Token, error)
}

type JwtTokenImpl struct{}

func NewJWT() *JwtTokenImpl {
	return &JwtTokenImpl{}
}

type JwtTokenClaims struct {
	Type   string
	Role   string
	UserID string
}

func NewJWTTokenClaims(userID, types, role string) *JwtTokenClaims {
	return &JwtTokenClaims{
		Type:   types,
		Role:   role,
		UserID: userID,
	}
}

func (j JwtTokenImpl) GenerateJwtTokenForAuth(types, userID, role string) (string, error) {

	type customClaims struct {
		Type string `json:"type"`
		Role string `json:"role"`
		jwt.RegisteredClaims
	}

	now := time.Now()

	registeredClaims := customClaims{
		Type: types,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   "mediseane",
			Subject:  userID,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(24 * time.Hour),
			},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, registeredClaims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", apperror.NewErrStatusUnauthorized(appconstant.FieldErrCheckAuthorization, apperror.ErrUnexpectedSigningMethod, nil)
	}
	return tokenString, nil
}

func (j JwtTokenImpl) ParseJwtTokenForAuth(c context.Context, tokenString string) (*JwtTokenClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperror.ErrTokenInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
		jwt.WithIssuedAt(),
		jwt.WithIssuer("mediseane"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	claims, exists := token.Claims.(jwt.MapClaims)
	if !exists {
		return nil, err
	}

	userID, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}

	types := claims["type"].(string)
	role := claims["role"].(string)
	jwtTokenClaims := NewJWTTokenClaims(userID, types, role)

	return jwtTokenClaims, nil
}
