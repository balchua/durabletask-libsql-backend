/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/balchua/durabletask-libsql-sample/pkg/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

type ServeCommand struct {
	serveCmd  *cobra.Command
	wfHandler *handler.WorkflowHandler
	port      int
}

func NewServeCommand(handler *handler.WorkflowHandler) *ServeCommand {
	serve := &ServeCommand{
		wfHandler: handler,
	}
	// serveCmd represents the serve command
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "starts a server",
		Long:  `Starts a server that listens to a port.`,
		Run:   serve.serve,
	}
	serve.serveCmd = serveCmd
	serve.setupFlags(serveCmd)
	return serve
}

func (s *ServeCommand) GetCommand() *cobra.Command {
	return s.serveCmd
}

func (s *ServeCommand) setupFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().IntVarP(&s.port, "port", "p", 8888, "port for the server to listen to")
}

func (s *ServeCommand) serve(cmd *cobra.Command, args []string) {
	app := fiber.New()
	app.Post("/workflow/simple", s.wfHandler.ScheduleSimpleOrchestration)
	app.Get("/status/orchestration/:id", s.wfHandler.GetStatus)

	if err := app.Listen(fmt.Sprintf(":%d", s.port)); err != nil {
		panic(err)
	}
}
