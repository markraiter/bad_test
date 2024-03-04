package handler

import (
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/markraiter/bad_test/internal/model"
)

type Processor interface {
	// FindValues finds min, max, median, average, max increasing sequential and max decreasing sequential of the numbers in the file.
	//
	// If file is not valid, returns error.
	// If file is empty, returns error.
	// If file is valid, returns min, max, median, average, max increasing sequential and max decreasing sequential of the numbers.
	FindValues(form *multipart.Form) (*model.TaskResult, error)
}

type TaskHandler struct {
	log       *slog.Logger
	processor Processor
}

// @Summary Find values
// @Description Web service that receives a .txt file with numbers and returns `min`, `max`, `median`, `average`, `max increasing sequential` and `max decreasing sequential` of the numbers.
// @Tags Task
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Please insert your .txt file here"
// @Success 200 {object} model.Response
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /task [post].
func (th *TaskHandler) FindValues(c *fiber.Ctx) error {
	const operation = "handler.Task.FindValues"

	log := th.log.With(slog.String("operation", operation))

	log.Info("processing request...")

	startTime := time.Now()

	form, err := c.MultipartForm()
	if err != nil {
		log.Error("error parsing form", model.Err(err))

		return c.Status(fiber.StatusBadRequest).JSON(model.Response{Message: err.Error()})
	}

	res, err := th.processor.FindValues(form)
	if err != nil {
		log.Error("error processing file", model.Err(err))

		return c.Status(fiber.StatusInternalServerError).JSON(model.Response{Message: err.Error()})
	}

	processingTime := time.Since(startTime)

	res.Time = processingTime.String()

	return c.Status(fiber.StatusOK).JSON(res)
}
