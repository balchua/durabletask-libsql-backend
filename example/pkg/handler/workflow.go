package handler

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"

	"github.com/balchua/durabletask-libsql-sample/pkg/orchestration"
	"github.com/gofiber/fiber/v2"
	"github.com/microsoft/durabletask-go/api"
)

type WorkflowHandler struct {
	orchestration *orchestration.Orchestration
}

func NewWorkflowHandler(o *orchestration.Orchestration) *WorkflowHandler {
	return &WorkflowHandler{
		orchestration: o,
	}
}

func (wf *WorkflowHandler) ScheduleSimpleOrchestration(c *fiber.Ctx) error {
	ctx := c.Context()
	id, err := wf.orchestration.ScheduleWorkflow(ctx, "SimpleOrchestration", api.WithInput(rand.Intn(100)))
	if err != nil {
		slog.Error("Failed to schedule new orchestration", slog.Any("error", err))
		return c.Status(http.StatusInternalServerError).SendString("Failed to schedule new orchestration")
	}

	slog.Info("Orchestration started", "id", id)
	c.Response().Header.Set("Location", fmt.Sprintf("/status/orchestration/%s", id))
	return c.Status(http.StatusAccepted).SendString("Orchestration started")
}

func (wf *WorkflowHandler) GetStatus(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")

	status, err := wf.orchestration.GetStatus(ctx, id)
	if err != nil {
		slog.Error("Failed to get status", slog.Any("error", err))
		return c.Status(http.StatusInternalServerError).SendString("Failed to get status")
	}
	return c.JSON(status)
}
