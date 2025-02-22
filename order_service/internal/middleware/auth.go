package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
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

			if claims["user_id"] == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Missing user_id claim",
				})
			}

			c.Locals("user_id", claims["user_id"])
			return c.Next()
		},
	})
}
