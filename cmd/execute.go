package cmd

import (
	"github.com/google/goterm/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"regexp"
	"remote_exec/util"
	"strings"
	"time"
)

var executeCmd = &cobra.Command{
	Use:     "execute",
	Short:   "execute command on remote",
	Long:    "execute command on remote",
	Example: "remote execute",
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
			port         = 22
			user         string
			passwd       string
			rootPwd      string
			hostList     []string
			hosts        []*util.Host
			commands     []string
			rootPrompt   *regexp.Regexp
			passwdPrompt *regexp.Regexp
			timeout      = int64(10)
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

		if config.IsSet("spc_hosts") {
			spcHosts := config.GetStringSlice("spc_hosts")
			for _, spcHost := range spcHosts {
				params := strings.Split(spcHost, " ")
				h := &util.Host{Port: port, User: user, Host: params[0], Passwd: params[1], RootPwd: params[2]}
				hosts = append(hosts, h)
			}
		}

		if config.IsSet("timeout") {
			timeout = config.GetInt64("timeout")
		}

		commands = command.GetStringSlice("commands")

		rootPrompt = regexp.MustCompile(config.GetString("rootPrompt"))
		passwdPrompt = regexp.MustCompile(config.GetString("passwdPrompt"))

		log.Println("start execute command...")
		for _, host := range hosts {
			for _, command := range commands {
				util.RemoteExec(host, command, rootPrompt, passwdPrompt, time.Duration(timeout)*time.Second)
			}
		}
		log.Println(term.Greenf("execute command finished."))
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
}
