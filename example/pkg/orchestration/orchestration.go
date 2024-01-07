package orchestration

import (
	"context"

	"github.com/microsoft/durabletask-go/api"
	"github.com/microsoft/durabletask-go/backend"
	"github.com/microsoft/durabletask-go/task"
)

type Orchestration struct {
	workflows     map[string]Workflow
	taskRegistry  *task.TaskRegistry
	be            backend.Backend
	logger        backend.Logger
	taskHubWorker backend.TaskHubWorker
	taskHubClient backend.TaskHubClient
}

func NewOrchestration(be backend.Backend, logger backend.Logger) *Orchestration {
	r := task.NewTaskRegistry()
	return &Orchestration{
		be:           be,
		taskRegistry: r,
		logger:       logger,
		workflows:    make(map[string]Workflow),
	}
}

func (o *Orchestration) AddWorkflow(workflowName string, w Workflow) {
	o.workflows[workflowName] = w
	o.registerWorkflows(w)
}

func (o *Orchestration) registerTasks(workflow Workflow) {
	tasks := workflow.GetTasks()
	for k, v := range tasks {
		o.taskRegistry.AddActivityN(k, v)
	}
}

func (o *Orchestration) registerWorkflows(w Workflow) {
	o.taskRegistry.AddOrchestrator(w.GetWorkflow())
	o.registerTasks(w)
}

func (o *Orchestration) initTaskHub(ctx context.Context) error {
	// Create an executor
	executor := task.NewTaskExecutor(o.taskRegistry)
	orchestrationWorker := backend.NewOrchestrationWorker(o.be, executor, o.logger)
	activityWorker := backend.NewActivityTaskWorker(o.be, executor, o.logger)
	o.taskHubWorker = backend.NewTaskHubWorker(o.be, orchestrationWorker, activityWorker, o.logger)

	// Start the worker
	err := o.taskHubWorker.Start(ctx)
	if err != nil {
		return err
	}

	// Get the client to the backend
	o.taskHubClient = backend.NewTaskHubClient(o.be)

	return nil
}

func (o *Orchestration) Start(ctx context.Context) error {
	if err := o.initTaskHub(ctx); err != nil {
		return err
	}
	return nil
}

func (o *Orchestration) Stop(ctx context.Context) error {
	return o.taskHubWorker.Shutdown(ctx)
}

func (o *Orchestration) ScheduleWorkflow(ctx context.Context, workflowId string, arg ...api.NewOrchestrationOptions) (api.InstanceID, error) {
	return o.taskHubClient.ScheduleNewOrchestration(ctx, workflowId, arg...)
}
