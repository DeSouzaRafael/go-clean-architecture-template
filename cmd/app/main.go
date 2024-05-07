package main

import (
	"log"

	"github.com/DeSouzaRafael/go-clean-architecture-template/config"
	_ "github.com/DeSouzaRafael/go-clean-architecture-template/docs"
	"github.com/DeSouzaRafael/go-clean-architecture-template/internal/app"
)

func main() {
	// configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// run app
	app.Run(cfg)
}
