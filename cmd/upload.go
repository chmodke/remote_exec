// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
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

		files = command.GetStringSlice("upload")

		log.Println("start upload file...")
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
				log.Printf("upload %s to %s %s.\n", from, host, to)
				util.RemotePut(22, host, user, passwd, from, to)
			}
		}
		log.Println("upload file finished.")
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
