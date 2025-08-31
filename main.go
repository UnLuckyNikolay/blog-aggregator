package main

import (
	"fmt"
	"log"
	"os"

	cmds "github.com/UnLuckyNikolay/blog-aggregator/internal/command_handlers"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/config"
)

func main() {
	var state cmds.State
	var cmdMan cmds.CommandManager
	var err error

	cmdMan.Initialize()

	state.Cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args // 0 - path, 1 - cmd name, 2+ - cmd args
	if len(args) < 2 {
		fmt.Println("No command found.")
		os.Exit(1)
	}
	err = cmdMan.HandleCommand(&state, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
