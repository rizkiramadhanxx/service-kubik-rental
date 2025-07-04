// middleware/authorization.go
package middleware

import (
	"encoding/json"
	"fmt"
	"kubik-rental/entity"

	"github.com/gofiber/fiber/v2"
)

// RequireModuleAccess ensures the user has access to a specific module
func RequireModuleAccess(module string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u, ok := c.Locals("user").(entity.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized or missing user context",
			})
		}

		var modules []string
		if err := json.Unmarshal([]byte(u.Role.Modules), &modules); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Failed to read role modules",
			})
		}

		for _, m := range modules {
			if m == module {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": fmt.Sprintf("Access denied to module: %s", module),
		})
	}
}
