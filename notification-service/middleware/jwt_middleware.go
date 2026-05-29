package middlewares

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const ContextKeyUserID = "user_id"

type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "missing authorization header",
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "invalid authorization header format",
				})
			}

			claims := &JWTClaims{}

			token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{
					"message": "invalid or expired token",
				})
			}

			c.Set(ContextKeyUserID, claims.UserID)

			return next(c)
		}
	}
}

func GetUserID(c echo.Context) string {
	userID, _ := c.Get(ContextKeyUserID).(string)
	return userID
}
