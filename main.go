package main

import (
	"log"
	"remote_exec/cmd"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime)
	cmd.Execute()
}
