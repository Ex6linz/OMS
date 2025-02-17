package middleware

package middleware

import (
    "github.com/gofiber/fiber/v2"
    jwtware "github.com/gofiber/jwt/v3"
    "github.com/golang-jwt/jwt/v4"
)

func NewAuthMiddleware(secret string) fiber.Handler {
    return jwtware.New(jwtware.Config{
        SigningKey:     []byte(secret),
        SigningMethod:  "HS256",
        TokenLookup:    "header:Authorization",
        AuthScheme:     "Bearer",
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Unauthorized",
                "message": err.Error(),
            })
        },
        SuccessHandler: func(c *fiber.Ctx) error {
            token := c.Locals("user").(*jwt.Token)
            claims := token.Claims.(jwt.MapClaims)
            
            // Verify required claims
            if claims["sub"] == nil {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                    "error": "Invalid token claims",
                })
            }
            
            return c.Next()
        },
    })
}