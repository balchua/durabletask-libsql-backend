package main

import (
	"context"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/balchua/durabletask-libsql-backend/lib/backend/libsql"
	"github.com/balchua/durabletask-libsql-sample/pkg/logging"
	"github.com/balchua/durabletask-libsql-sample/pkg/orchestration"
	"github.com/balchua/durabletask-libsql-sample/pkg/tracing"
	"github.com/balchua/durabletask-libsql-sample/pkg/workflow/simple"
	"github.com/microsoft/durabletask-go/api"
	"github.com/microsoft/durabletask-go/backend"
)

func main() {
	// initialize the logger
	logger := logging.NewSlogWrapper(logging.WithLogLevel(slog.LevelInfo))
	// init the tracing
	// Tracing can be configured independently of the orchestration code.
	tp, err := tracing.ConfigureZipkinTracing()
	if err != nil {
		logger.ErrorS("Failed to create tracer", "tracing error", err)
		panic(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logger.ErrorS("Failed to stop tracer", "tracing error", err)
			panic(err)
		}
	}()
	logger.InfoS("starting durable tasks", "orchestrator", "SimpleOrchestration")

	// Init the client
	ctx := context.Background()
	// logger := backend.DefaultLogger()
	be := getBackend(logger)
	workflows, err := InitOrchestrations(ctx, be, logger)
	if err != nil {
		logger.ErrorS("Failed to initialize the orchestrations", "error", err)
		panic(err)
	}
	defer workflows.Stop(ctx)

	// Start a new orchestration
	id, err := workflows.ScheduleWorkflow(ctx, "SimpleOrchestration", api.WithInput(rand.Intn(100)))
	if err != nil {
		logger.ErrorS("Failed to schedule new orchestration", "error", err)
		panic(err)
	}

	slog.Info("Orchestration started", "id", id)

	// // Wait for the orchestration to complete
	// metadata, err := client.WaitForOrchestrationCompletion(ctx, id)
	// if err != nil {
	// 	logger.ErrorS("Failed to wait for orchestration to complete", "error", err)
	// 	panic(err)
	// }
	// var emps []employee
	// json.Unmarshal([]byte(metadata.SerializedOutput), &emps)
	// // Print the results
	// logger.Infof("Orchestration completed: %v", emps)
	// // Cleanup the task hub
	// if err := be.DeleteTaskHub(ctx); err != nil {
	// 	logger.ErrorS("Failed to delete task hub: %v", err)
	// 	panic(err)
	// }

	end := time.After(60 * time.Second)
	<-end

}

func getBackend(logger backend.Logger) backend.Backend {
	token := os.Getenv("DURABLETASK_DEMO_DB_TOKEN")
	host := os.Getenv("DURABLETASK_DEMO_DB_HOST")
	scheme := os.Getenv("DURABLETASK_DEMO_DB_HOST_SCHEME")

	// Use the libsql-server storage provider
	return libsql.NewLibSqlBackend(libsql.WithScheme(scheme), libsql.WithHost(host), libsql.WithLogger(logger), libsql.WithToken(token))
}

func InitOrchestrations(ctx context.Context, be backend.Backend, logger backend.Logger) (*orchestration.Orchestration, error) {
	orchestration := orchestration.NewOrchestration(be, logger)
	simpleWorkflow := simple.NewSimpleWorkflow()
	orchestration.AddWorkflow("SimpleOrchestration", simpleWorkflow)
	orchestration.Start(ctx)
	return orchestration, nil

}
