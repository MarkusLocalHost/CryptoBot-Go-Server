package middleware

import (
	"cryptocurrency/internal/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationType       = "Bearer"
	AuthorizationPayloadKey = "data"
)

func AuthMiddleware(s models.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)
		if authorizationHeader == "" {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) != 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		accessToken := fields[1]
		accessToken = strings.ReplaceAll(accessToken, "b'", "")
		accessToken = strings.ReplaceAll(accessToken, "'", "")
		payload, err := s.ValidateIDToken(accessToken)
		if err != nil {
			err := errors.New("invalid token")
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
