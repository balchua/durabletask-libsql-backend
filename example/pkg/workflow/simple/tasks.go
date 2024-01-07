package simple

import (
	"fmt"
	"log/slog"

	"github.com/microsoft/durabletask-go/task"
)

// GetEmployeeDetailById can be called by an orchestrator function and will return a friendly greeting.
func (s *SimpleWorkflow) getEmployeeDetailById(ctx task.ActivityContext) (any, error) {
	var id string
	if err := ctx.GetInput(&id); err != nil {
		return nil, err
	}
	slog.Info("employee id", "value", id)

	for i, e := range employees {
		if e.Id == id {
			return &employees[i], nil
		}
	}
	return nil, fmt.Errorf("employee not found")

}

func (s *SimpleWorkflow) task2(ctx task.ActivityContext) (any, error) {
	slog.Info("executing task2")
	return "not doing anything", nil
}
