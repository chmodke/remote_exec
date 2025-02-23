package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"regexp"
	"remote_exec/goterm/term"
	"remote_exec/util"
	"time"
)

var execCmd = &cobra.Command{
	Use:     "exec [flags] [section]",
	Short:   "execute command on remote",
	Long:    "execute command on remote",
	Example: "remote exec",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err            error
			config         *viper.Viper
			command        *viper.Viper
			hosts          []*util.Host
			commands       []string
			rootPrompt     *regexp.Regexp
			passwdPrompt   *regexp.Regexp
			timeout        = int64(10)
			thread, _      = cmd.Flags().GetInt(util.ConstThread)
			configPath, _  = cmd.Flags().GetString(util.ConstConfig)
			commandPath, _ = cmd.Flags().GetString(util.ConstCommand)
		)

		if hosts, err = util.ParseHosts(configPath, cmd); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}
		if config, err = util.LoadCfg(configPath, util.DefaultConfig); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}
		if command, err = util.LoadCfg(commandPath, util.DefaultCommand); err != nil {
			log.Println(term.Redf(err.Error()))
			return
		}

		if config.IsSet("timeout") {
			timeout = config.GetInt64("timeout")
		}

		if len(args) > 0 {
			section := args[0]
			if !command.InConfig(section) {
				log.Println(term.Redf("no %s configuration item found.", section))
				return
			}
			commands = command.Sub(section).GetStringSlice("exec")
		} else {
			commands = command.GetStringSlice("exec")
		}

		rootPrompt = regexp.MustCompile(config.GetString("root-prompt"))
		passwdPrompt = regexp.MustCompile(config.GetString("passwd-prompt"))

		log.Println("start execute command...")

		util.Process(thread, hosts, func(host *util.Host) {
			util.RemoteExec(host, commands, rootPrompt, passwdPrompt, time.Duration(timeout)*time.Second)

		})
		log.Println(term.Greenf("execute command finished."))

	},
}

func init() {
	rootCmd.AddCommand(execCmd)
}
