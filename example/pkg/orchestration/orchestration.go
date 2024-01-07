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

func (o *Orchestration) AddWorkflow(workflowName string, w Workflow) error {
	o.workflows[workflowName] = w
	return o.registerWorkflows(workflowName, w)
}

func (o *Orchestration) registerTasks(workflow Workflow) error {
	tasks := workflow.GetTasks()
	for k, v := range tasks {
		if err := o.taskRegistry.AddActivityN(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (o *Orchestration) registerWorkflows(workflowName string, w Workflow) error {
	if err := o.taskRegistry.AddOrchestratorN(workflowName, w.GetWorkflow()); err != nil {
		return err
	}
	return o.registerTasks(w)
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

func (o *Orchestration) isValidWorkflow(workflowName string) bool {
	if _, ok := o.workflows[workflowName]; !ok {
		return false
	}
	return true
}

func (o *Orchestration) ScheduleWorkflow(ctx context.Context, workflowName string, arg ...api.NewOrchestrationOptions) (api.InstanceID, error) {
	if o.isValidWorkflow(workflowName) {
		return o.taskHubClient.ScheduleNewOrchestration(ctx, workflowName, arg...)
	}
	return "", ErrWorkflowNotFound
}

func (o *Orchestration) GetStatus(ctx context.Context, instanceID string) (string, error) {
	var (
		metadata *api.OrchestrationMetadata
		err      error
	)

	if metadata, err = o.taskHubClient.FetchOrchestrationMetadata(ctx, api.InstanceID(instanceID)); err != nil {
		return api.RUNTIME_STATUS_FAILED.String(), err
	}
	return metadata.RuntimeStatus.String(), nil
}
