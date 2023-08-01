package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"regexp"
	"remote_exec/util"
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
			log.Fatalln("load config.yaml failed.")
		}

		if command, err = util.LoadCfg("command"); err != nil {
			log.Fatalln("load command.yaml failed.")
		}
		var (
			user         string
			passwd       string
			rootPwd      string
			hosts        []string
			commands     []string
			rootPrompt   *regexp.Regexp
			passwdPrompt *regexp.Regexp
			timeout      int64
		)
		user = config.GetString("user")
		passwd = config.GetString("passwd")
		rootPwd = config.GetString("rootPwd")
		hosts = config.GetStringSlice("hosts")
		timeout = config.GetInt64("timeout")

		commands = command.GetStringSlice("commands")

		rootPrompt = regexp.MustCompile(config.GetString("rootPrompt"))
		passwdPrompt = regexp.MustCompile(config.GetString("passwdPrompt"))

		for _, host := range hosts {
			for _, command := range commands {
				log.Printf("execute %s on %s.\n", command, host)
				util.RemoteExec(22, host, user, passwd, rootPwd, command, rootPrompt, passwdPrompt, timeout)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
}
