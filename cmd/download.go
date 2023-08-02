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

var downloadCmd = &cobra.Command{
	Use:     "download",
	Short:   "download file to remote",
	Long:    "download file to remote",
	Example: "remote download",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			config  *viper.Viper
			command *viper.Viper
			err     error
		)
		if config, err = util.LoadCfg("config"); err != nil {
			log.Println(term.Redf("load config.yaml failed."))
			return
		}
		if command, err = util.LoadCfg("command"); err != nil {
			log.Println(term.Redf("load command.yaml failed."))
			return
		}
		var (
			user   string
			passwd string
			hosts  []string
			files  []string
		)
		user = config.GetString("user")
		passwd = config.GetString("passwd")
		hosts = config.GetStringSlice("hosts")

		files = command.GetStringSlice("download")
		log.Println("start download file...")
		for _, host := range hosts {
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
				util.RemoteGet(22, host, user, passwd, from, to)
			}
		}
		log.Println(term.Greenf("download file finished."))
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
