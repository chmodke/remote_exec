// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"remote_exec/goterm/term"
	"remote_exec/util"
)

var putCmd = &cobra.Command{
	Use:     "put [flags] [section]",
	Short:   "put file to remote",
	Long:    "put file to remote",
	Example: "remote put",
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
			files = command.Sub(section).GetStringSlice("put")
		} else {
			files = command.GetStringSlice("put")
		}

		log.Println("start put file...")

		util.Process(thread, hosts, func(host *util.Host) {
			util.RemotePut(host, files)
		})
		log.Println(term.Greenf("put file finished."))
	},
}

func init() {
	rootCmd.AddCommand(putCmd)
}
