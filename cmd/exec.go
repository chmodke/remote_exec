package cmd

import (
	"github.com/google/goterm/term"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"regexp"
	"remote_exec/util"
	"time"
)

var execCmd = &cobra.Command{
	Use:     "exec",
	Short:   "execute command on remote",
	Long:    "execute command on remote",
	Example: "remote exec",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err          error
			config       *viper.Viper
			command      *viper.Viper
			hosts        []*util.Host
			commands     []string
			rootPrompt   *regexp.Regexp
			passwdPrompt *regexp.Regexp
			timeout      = int64(10)
			thread       int
		)

		if hosts, err = util.ParseHosts(); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}
		if config, err = util.LoadCfg("config"); err != nil {
			log.Println(term.Redf("load config.yaml failed."))
			return
		}
		if command, err = util.LoadCfg("command"); err != nil {
			log.Println(term.Redf("load command.yaml failed."))
			return
		}

		if config.IsSet("timeout") {
			timeout = config.GetInt64("timeout")
		}

		commands = command.GetStringSlice("exec")

		rootPrompt = regexp.MustCompile(config.GetString("rootPrompt"))
		passwdPrompt = regexp.MustCompile(config.GetString("passwdPrompt"))

		log.Println("start execute command...")

		if thread, err = cmd.Flags().GetInt("thread"); err != nil {
			thread = 1
		}
		util.Process(thread, hosts, commands, func(host *util.Host, commands []string) {
			for _, command := range commands {
				util.RemoteExec(host, command, rootPrompt, passwdPrompt, time.Duration(timeout)*time.Second)
			}
		})
		log.Println(term.Greenf("execute command finished."))

	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
