// Package cmd /*
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"remote_exec/util"
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
	rootCmd.Version = fmt.Sprintf("%s\nReport Bug: <https://gitee.com/chmodke/remote_exec/>", util.VERSION)
	rootCmd.PersistentFlags().IntP(util.ConstThread, "t", util.DefaultThread, "maximum number of concurrent (0 < t <= 16)")
	rootCmd.PersistentFlags().StringP(util.ConstConfig, "f", util.DefaultConfig, "Specify servers configuration")
	rootCmd.PersistentFlags().StringP(util.ConstCommand, "c", util.DefaultCommand, "Specify commands configuration")
	rootCmd.PersistentFlags().StringP(util.ConstNetMask, "m", "", "ip filter, e.g. 192.168.1.1 192.168.1.1,192.168.1.2 192.168.0.0/24")
}
