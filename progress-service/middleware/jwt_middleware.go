package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	ContextKeyUserID = "user_id"
	ContextKeyRole   = "role"

	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleAdmin   = "admin"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": "missing authorization header",
					"data":    nil,
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": "invalid authorization header format",
					"data":    nil,
				})
			}

			tokenStr := parts[1]
			claims := &JWTClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.ErrUnauthorized
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"success": false,
					"message": "invalid or expired token",
					"data":    nil,
				})
			}

			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyRole, claims.Role)
			return next(c)
		}
	}
}

// RequireRole returns middleware that enforces one of the allowed roles.
func RequireRole(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get(ContextKeyRole).(string)
			if !ok || !allowed[role] {
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"success": false,
					"message": "forbidden: insufficient role",
					"data":    nil,
				})
			}
			return next(c)
		}
	}
}

// GetUserID extracts the user_id from the echo context.
func GetUserID(c echo.Context) string {
	v, _ := c.Get(ContextKeyUserID).(string)
	return v
}

// GetRole extracts the role from the echo context.
func GetRole(c echo.Context) string {
	v, _ := c.Get(ContextKeyRole).(string)
	return v
}
