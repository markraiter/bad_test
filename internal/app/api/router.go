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

	app.Post(apiPrefix+"/task", timeout.NewWithContext(handler.TaskHandler.FindValues, cfg.Server.AppWriteTimeout))
}
