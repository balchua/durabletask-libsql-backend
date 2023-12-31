package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/balchua/durabletask-libsql-sample/pkg/logging"
	"github.com/balchua/durabletask-libsql/backend/libsql"
	"github.com/microsoft/durabletask-go/api"
	"github.com/microsoft/durabletask-go/backend"
	"github.com/microsoft/durabletask-go/task"
)

type employee struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

var employees = []employee{
	{"1", "John", "Doe", 30},
	{"2", "Jane", "Doe", 25},
	{"3", "Steven", "Smith", 40},
	{"4", "Lily", "Smith", 35},
}

func main() {
	// initialize the logger
	logger := logging.NewSlogWrapper(logging.WithLogLevel(slog.LevelInfo))
	logger.InfoS("starting durable tasks", "orchestrator", "SimpleOrchestration")
	// Create a new task registry and add the orchestrator and activities
	r := task.NewTaskRegistry()
	r.AddOrchestrator(SimpleOrchestration)
	r.AddActivity(GetEmployeeDetailById)
	// Init the client
	ctx := context.Background()
	// logger := backend.DefaultLogger()
	be := getBackend(logger)
	client, worker, err := Init(ctx, r, be, logger)
	if err != nil {
		logger.ErrorS("Failed to initialize the client", "error", err)
		panic(err)
	}
	defer worker.Shutdown(ctx)

	// Start a new orchestration
	id, err := client.ScheduleNewOrchestration(ctx, SimpleOrchestration, api.WithInput(rand.Intn(100)))
	if err != nil {
		logger.ErrorS("Failed to schedule new orchestration", "error", err)
		panic(err)
	}

	// Wait for the orchestration to complete
	metadata, err := client.WaitForOrchestrationCompletion(ctx, id)
	if err != nil {
		logger.ErrorS("Failed to wait for orchestration to complete", "error", err)
		panic(err)
	}
	var emps []employee
	json.Unmarshal([]byte(metadata.SerializedOutput), &emps)
	// Print the results
	logger.Infof("Orchestration completed: %v", emps)
	// // Cleanup the task hub
	// if err := be.DeleteTaskHub(ctx); err != nil {
	// 	logger.ErrorS("Failed to delete task hub: %v", err)
	// 	panic(err)
	// }

}

func getBackend(logger backend.Logger) backend.Backend {
	token := os.Getenv("DURABLETASK_DEMO_DB_TOKEN")
	host := os.Getenv("DURABLETASK_DEMO_DB_HOST")
	scheme := os.Getenv("DURABLETASK_DEMO_DB_HOST_SCHEME")
	// Use the libsql-server storage provider
	return libsql.NewLibSqlBackend(libsql.NewLibSqlOptions(scheme, host, token, 20*time.Second, 20*time.Second), logger)
}

// Init creates and initializes an in-memory client and worker pair with default configuration.
func Init(ctx context.Context, r *task.TaskRegistry, be backend.Backend, logger backend.Logger) (backend.TaskHubClient, backend.TaskHubWorker, error) {

	// Create an executor
	executor := task.NewTaskExecutor(r)
	orchestrationWorker := backend.NewOrchestrationWorker(be, executor, logger)
	activityWorker := backend.NewActivityTaskWorker(be, executor, logger)
	taskHubWorker := backend.NewTaskHubWorker(be, orchestrationWorker, activityWorker, logger)

	// Start the worker
	err := taskHubWorker.Start(ctx)
	if err != nil {
		return nil, nil, err
	}

	// Get the client to the backend
	taskHubClient := backend.NewTaskHubClient(be)

	return taskHubClient, taskHubWorker, nil
}

// SimpleOrchestration makes 2 activity calls in sequence and results the results
// as an array.
func SimpleOrchestration(ctx *task.OrchestrationContext) (any, error) {
	var input int
	ctx.GetInput(&input)
	slog.Info("input", "value", input)

	var john employee
	if err := ctx.CallActivity(GetEmployeeDetailById, task.WithActivityInput("1")).Await(&john); err != nil {
		return nil, err
	}
	var lily employee
	if err := ctx.CallActivity(GetEmployeeDetailById, task.WithActivityInput("4")).Await(&lily); err != nil {
		return nil, err
	}

	var jane employee
	if err := ctx.CallActivity(GetEmployeeDetailById, task.WithActivityInput("2")).Await(&jane); err != nil {
		return nil, err
	}

	var steven employee
	if err := ctx.CallActivity(GetEmployeeDetailById, task.WithActivityInput("2")).Await(&steven); err != nil {
		return nil, err
	}

	return []employee{john, lily, jane, steven}, nil
}

// GetEmployeeDetailById can be called by an orchestrator function and will return a friendly greeting.
func GetEmployeeDetailById(ctx task.ActivityContext) (any, error) {
	var id string
	if err := ctx.GetInput(&id); err != nil {
		return nil, err
	}

	for i, e := range employees {
		if e.Id == id {
			return &employees[i], nil
		}
	}
	return nil, fmt.Errorf("employee not found")

}
