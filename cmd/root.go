// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "remote",
	Short: "remote execute tool",
	Long:  "remote execute tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = "1.0.0\nReport Bug: <https://gitee.com/chmodke/remote_exec/>"
	rootCmd.PersistentFlags().IntP("thread", "t", 1, "maximum number of concurrent (0 < t <= 16)")
}
