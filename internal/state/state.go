package state

import (
	"net/http"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/config"
	"github.com/UnLuckyNikolay/blog-aggregator/internal/database"
)

type State struct {
	Cfg        *config.Config
	Db         *database.Queries
	HttpClient *http.Client
}
