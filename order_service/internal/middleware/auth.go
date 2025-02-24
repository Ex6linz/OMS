package middleware

import (
	"github.com/Ex6linz/OMS/order-service/internal/logger"
	"github.com/Ex6linz/OMS/order-service/internal/rbac"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:  []byte(secret),
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": err.Error(),
			})
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			tokenVal := c.Locals("user")
			token, ok := tokenVal.(*jwt.Token)
			if !ok || token == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token format",
				})
			}

			if !token.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token",
				})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token claims format",
				})
			}

			if claims["user"] == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing user claim",
				})
			}

			role, ok := claims["role"].(string)
			if !ok {
				logger.Log.Warn("No role in token", zap.Any("claims", claims))
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"errror": "Forbidenn missing user role",
				})
			}
			c.Locals("role", role)

			perms := rbac.GetPermissions(role)
			c.Locals("permissions", perms)

			c.Locals("user", claims["user"])
			return c.Next()
		},
	})
}
