// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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
			log.Fatalln("load config.yaml failed.")
		}
		if command, err = util.LoadCfg("command"); err != nil {
			log.Fatalln("load command.yaml failed.")
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
				to = params[0]
				if len(params) == 2 {
					to = params[1]
				}
				log.Printf("from %s download %s write to %s.\n", host, from, to)
				util.RemoteGet(22, host, user, passwd, from, to)
			}
		}
		log.Println("download file finished.")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
