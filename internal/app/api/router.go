package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/gofiber/swagger"
	"github.com/markraiter/bad_test/internal/app/api/handler"
	"github.com/markraiter/bad_test/internal/config"
)

const apiPrefix = "/api/v1"

// initRoutes configures the routes for the app.
func (s Server) initRoutes(app *fiber.App, handler *handler.Handler, cfg *config.Config) {
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Get(apiPrefix+"/health", timeout.NewWithContext(handler.APIHealth, cfg.Server.AppReadTimeout))

	// app.Post(apiPrefix+"/values", timeout.NewWithContext(handler.Values, cfg.Server.AppReadTimeout))
}
