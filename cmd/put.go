// Package cmd /*
package cmd

import (
	"github.com/google/goterm/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path"
	"remote_exec/util"
	"strings"
)

var putCmd = &cobra.Command{
	Use:     "put",
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
		if hosts, err = util.ParseHosts(configPath); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}
		if command, err = util.LoadCfg(commandPath, util.DefaultCommand); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}

		files = command.GetStringSlice("put")

		log.Println("start put file...")

		util.Process(thread, hosts, files, func(host *util.Host, files []string) {
			for _, file := range files {
				params := strings.Split(file, "#")
				var (
					from string
					to   string
				)
				from = params[0]
				to = path.Dir(params[0])
				if len(params) == 2 {
					to = params[1]
				}
				util.RemotePut(host, from, to)
			}
		})
		log.Println(term.Greenf("put file finished."))
	},
}

func init() {
	rootCmd.AddCommand(putCmd)
}
