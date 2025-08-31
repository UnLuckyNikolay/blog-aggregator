package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	cmds "github.com/UnLuckyNikolay/blog-aggregator/internal/command_handlers"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/config"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/database"
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

	db, err := sql.Open("postgres", state.Cfg.DbUrl)
	state.Db = database.New(db)

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
