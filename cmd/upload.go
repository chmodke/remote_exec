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

var uploadCmd = &cobra.Command{
	Use:     "upload",
	Short:   "upload file to remote",
	Long:    "upload file to remote",
	Example: "remote upload",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err     error
			command *viper.Viper
			hosts   []*util.Host
			files   []string
			thread  int
		)
		if hosts, err = util.ParseHosts(); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}
		if command, err = util.LoadCfg("command"); err != nil {
			log.Println(term.Redf("load command.yaml failed."))
			return
		}

		files = command.GetStringSlice("upload")

		log.Println("start upload file...")

		if thread, err = cmd.Flags().GetInt("thread"); err != nil {
			thread = 1
		}
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
		log.Println(term.Greenf("upload file finished."))
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
