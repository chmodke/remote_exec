package main

import (
	"remote_exec/cmd"
	"remote_exec/util"
)

func main() {
	if closeFile := util.InitLog(); closeFile != nil {
		defer closeFile()
	}
	cmd.Execute()

}
