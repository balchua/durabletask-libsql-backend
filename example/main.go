/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/balchua/durabletask-libsql-backend/lib/backend/libsql"
	"github.com/balchua/durabletask-libsql-sample/cmd"
	"github.com/balchua/durabletask-libsql-sample/pkg/handler"
	"github.com/balchua/durabletask-libsql-sample/pkg/orchestration"
	"github.com/balchua/durabletask-libsql-sample/pkg/workflow/simple"
	"github.com/microsoft/durabletask-go/backend"
)

func main() {

	rootCmd := cmd.NewRootCommand()
	logger := rootCmd.Logger()
	be := initBackend(logger)
	orchestration, err := initOrchestration(context.Background(), be, logger)
	wf := handler.NewWorkflowHandler(orchestration)
	if err != nil {
		slog.Error("Failed to initialize orchestration", slog.Any("error", err))
		os.Exit(1)
	}
	serveCmd := cmd.NewServeCommand(wf)

	rootCmd.AddCommand(serveCmd.GetCommand())
	rootCmd.Execute()
}

func initOrchestration(ctx context.Context, be backend.Backend, logger backend.Logger) (*orchestration.Orchestration, error) {
	orchestration := orchestration.NewOrchestration(be, logger)
	simpleWorkflow := simple.NewSimpleWorkflow()
	if err := orchestration.AddWorkflow("SimpleOrchestration", simpleWorkflow); err != nil {
		slog.Error("Failed to add workflow", slog.Any("error", err))
		return nil, err
	}
	orchestration.Start(ctx)
	return orchestration, nil

}

func initBackend(logger backend.Logger) backend.Backend {
	token := os.Getenv("DURABLETASK_DEMO_DB_TOKEN")
	host := os.Getenv("DURABLETASK_DEMO_DB_HOST")
	scheme := os.Getenv("DURABLETASK_DEMO_DB_HOST_SCHEME")

	// Use the libsql-server storage provider
	return libsql.NewLibSqlBackend(libsql.WithScheme(scheme), libsql.WithHost(host), libsql.WithLogger(logger), libsql.WithToken(token))
}
