package main

import (
	"btcwallet/internal/app/api"
	"flag"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.json", "path to config")
}

func main() {
	flag.Parse()

	config := api.NewConfig()
	if err := config.ReadConfig(configPath); err != nil {
		log.Fatalf("Error reading configs file: %s\n", err)
	}

	log.Fatal(api.Start(config))
}
