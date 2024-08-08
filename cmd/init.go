package cmd

import (
	"github.com/spf13/cobra"
	"grpc-template-service/cmd/config"
	"grpc-template-service/cmd/create"
	"grpc-template-service/cmd/server"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          "mod",
	Short:        "mod",
	SilenceUsage: true,
	Long:         `mod`,
}

func init() {
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(create.StartCmd)
	rootCmd.AddCommand(server.StartCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
