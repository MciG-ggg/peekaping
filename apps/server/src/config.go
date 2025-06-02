package main

import (
	"encoding/json"
	"peekaping/src/config"
)

func ProvideConfig() (*config.Config, error) {
	cfg, err := config.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	// Convert struct to JSON and print it
	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		println("Error marshaling config to JSON:", "error", err.Error())
		panic(err)
	} else {
		println("Configuration loaded:")
		println(string(configJSON))
	}

	return &cfg, nil
}
