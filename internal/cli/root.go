package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "echostrike",
	Short: "EchoStrike: A high-performance syslog generator and attack simulator",
	Long: `EchoStrike is a CLI tool designed for security professionals to generate
realistic syslog data, simulate attacks and test SIEM pipelines.

It supports TCP, UDP, TLS, custom templates and high-volume traffic generation.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do generic stuff or help
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
