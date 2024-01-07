/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log/slog"
	"os"

	"github.com/balchua/durabletask-libsql-sample/pkg/logging"
	"github.com/microsoft/durabletask-go/backend"
	"github.com/spf13/cobra"
)

var debug bool

type RootCommand struct {
	rootCmd *cobra.Command
	logger  logging.SlogWrapper
}

func NewRootCommand() *RootCommand {
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "durabletask-libsql-sample",
		Short: "This is a demo of using durabletask with libsql backend",
		Long:  `This is a demo of using durabletask with libsql backend.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}
	setFlags(rootCmd)
	return &RootCommand{
		rootCmd: rootCmd,
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func (r *RootCommand) Execute() {
	err := r.rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func (r *RootCommand) AddCommand(cmd *cobra.Command) {
	r.rootCmd.AddCommand(cmd)
}

func (r *RootCommand) Logger() backend.Logger {
	r.rootCmd.Flags()
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}
	r.logger = logging.NewSlogWrapper(logging.WithLogLevel(logLevel))
	return r.logger
}

func setFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
}
