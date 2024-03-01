package handler

import (
	"log/slog"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/markraiter/bad_test/internal/model"
)

type Processor interface {
	FindValues(form *multipart.Form) (string, error)
}

type TaskHandler struct {
	log       *slog.Logger
	processor Processor
}

// @Summary Find values
// @Description Find values
// @Tags Task
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /task [post].
func (th *TaskHandler) FindValues(c *fiber.Ctx) error {
	const operation = "handler.Task.FindValues"

	log := th.log.With(slog.String("operation", operation))

	log.Info("attempting to process file")

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: err.Error()})
	}

	res, err := th.processor.FindValues(form)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{Message: err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(model.Response{Message: res})
}
