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
			port     = 22
			user     string
			passwd   string
			rootPwd  string
			hostList []string
			hosts    []*util.Host
			files    []string
		)
		if config.IsSet("port") {
			port = config.GetInt("port")
		}
		user = config.GetString("user")
		passwd = config.GetString("passwd")
		rootPwd = config.GetString("rootPwd")
		hostList = config.GetStringSlice("hosts")

		for _, host := range hostList {
			h := &util.Host{Port: port, Host: host, User: user, Passwd: passwd, RootPwd: rootPwd}
			hosts = append(hosts, h)
		}

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
				to = path.Dir(params[0])
				if len(params) == 2 {
					to = params[1]
				}
				util.RemotePut(host, from, to)
			}
		}
		log.Println(term.Greenf("upload file finished."))
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
