package main

import (
	"github.com/tonradar/ton-dice-web-resolver/config"
	"github.com/tonradar/ton-dice-web-resolver/resolver"
)

func main() {
	cfg := config.GetConfig()

	service := resolver.NewResolver(&cfg)
	service.Start()
}
