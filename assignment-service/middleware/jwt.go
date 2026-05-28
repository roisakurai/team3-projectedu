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

	RoleTeacher = "teacher"
	RoleStudent = "student"
	RoleAdmin   = "admin"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTMiddleware validates the Bearer token and injects user_id + role into context.
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return respondUnauthorized(c, "missing authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				return respondUnauthorized(c, "invalid authorization header format")
			}

			tokenStr := parts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				return respondUnauthorized(c, "invalid or expired token")
			}

			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyRole, claims.Role)

			return next(c)
		}
	}
}

// RequireRole returns middleware that enforces one of the provided roles.
func RequireRole(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, ok := c.Get(ContextKeyRole).(string)
			if !ok || role == "" {
				return respondUnauthorized(c, "role not found in context")
			}
			if _, permitted := allowed[role]; !permitted {
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

func respondUnauthorized(c echo.Context, msg string) error {
	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"success": false,
		"message": msg,
		"data":    nil,
	})
}
