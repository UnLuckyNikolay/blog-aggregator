package main

import (
	"fmt"
	"log"

	"github.com/UnLuckyNikolay/blog-aggregator/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	cfg.SetUser("Niko")

	cfg, err = config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg)
}
