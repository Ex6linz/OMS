package utils

import (
	"github.com/gofiber/fiber/v2"
)

func HasPermission(c *fiber.Ctx, permission string) bool {
	perms, ok := c.Locals("permissions").([]string)
	if !ok {
		return false
	}
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}
