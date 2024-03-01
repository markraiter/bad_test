package handler

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/markraiter/bad_test/internal/model"
)

type HealthCheck struct {
	log *slog.Logger
}

// @Summary	Shows the status of server.
// @Description	Ping health of API for Docker.
// @Tags Health
// @Accept */*
// @Produce json
// @Success	200	{object} model.Response
// @Router /health [get].
func (hc *HealthCheck) APIHealth(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).JSON(model.Response{Message: "healthy"})
}
