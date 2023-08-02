package cmd

import (
	"github.com/google/goterm/term"
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
			log.Println(term.Redf("load config.yaml failed."))
			return
		}
		if command, err = util.LoadCfg("command"); err != nil {
			log.Println(term.Redf("load command.yaml failed."))
			return
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

		log.Println("start execute command...")
		for _, host := range hosts {
			for _, command := range commands {
				util.RemoteExec(22, host, user, passwd, rootPwd, command, rootPrompt, passwdPrompt, timeout)
			}
		}
		log.Println(term.Greenf("execute command finished."))
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
}
