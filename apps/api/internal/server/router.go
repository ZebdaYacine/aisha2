package httpapi

import "github.com/gofiber/fiber/v3"

// APIGroup is the single public API mount point. Feature route packages are
// mounted below this group by the server composition layer.
func APIGroup(app *fiber.App) fiber.Router {
	return app.Group("/api/v1")
}
