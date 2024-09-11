// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"remote_exec/goterm/term"
	"remote_exec/util"
)

var getCmd = &cobra.Command{
	Use:     "get [flags] [section]",
	Short:   "get file from remote",
	Long:    "get file from remote",
	Example: "remote get",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err            error
			command        *viper.Viper
			hosts          []*util.Host
			files          []string
			thread, _      = cmd.Flags().GetInt(util.ConstThread)
			configPath, _  = cmd.Flags().GetString(util.ConstConfig)
			commandPath, _ = cmd.Flags().GetString(util.ConstCommand)
		)
		if hosts, err = util.ParseHosts(configPath, cmd); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}
		if command, err = util.LoadCfg(commandPath, util.DefaultCommand); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}

		if len(args) > 0 {
			section := args[0]
			if !command.InConfig(section) {
				log.Println(term.Redf("no %s configuration item found.", section))
				return
			}
			files = command.Sub(section).GetStringSlice("get")
		} else {
			files = command.GetStringSlice("get")
		}

		log.Println("start get file...")

		util.Process(thread, hosts, func(host *util.Host) {
			util.RemoteGet(host, files)
		})
		log.Println(term.Greenf("get file finished."))
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
