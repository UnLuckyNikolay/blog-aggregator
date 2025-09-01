package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/commands"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/config"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/database"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/rssfeed"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/state"
)

func main() {
	var state state.State
	var cmdMan commands.CommandManager
	var err error

	cmdMan.Initialize()
	state.HttpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
	state.Cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", state.Cfg.DbUrl)
	state.Db = database.New(db)

	/*args := os.Args // 0 - path, 1 - cmd name, 2+ - cmd args
	if len(args) < 2 {
		fmt.Println("No command found.")
		os.Exit(1)
	}
	err = cmdMan.HandleCommand(&state, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/

	rssfeed.Agg(&state)
}
