package orchestration

import "github.com/microsoft/durabletask-go/task"

type Workflow interface {
	GetTasks() map[string]func(ctx task.ActivityContext) (any, error)
	GetWorkflow() func(ctx *task.OrchestrationContext) (any, error)
}
