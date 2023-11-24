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

var executeCmd = &cobra.Command{
	Use:     "execute",
	Short:   "execute command on remote",
	Long:    "execute command on remote",
	Example: "remote execute",
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
			errCnt       = 0
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

		commands = command.GetStringSlice("commands")

		rootPrompt = regexp.MustCompile(config.GetString("rootPrompt"))
		passwdPrompt = regexp.MustCompile(config.GetString("passwdPrompt"))

		log.Println("start execute command...")
		for idx, host := range hosts {
			result := true
			for _, command := range commands {
				result = util.RemoteExec(host, command, rootPrompt, passwdPrompt, time.Duration(timeout)*time.Second) && result
			}
			if !result {
				errCnt++
			}
			log.Printf("progress [%v/%v/%v]...\n", idx+1, errCnt, len(hosts))
		}
		log.Println(term.Greenf("execute command finished."))
	},
}

func init() {
	rootCmd.AddCommand(executeCmd)
}
