package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Symbols []string `json:"symbols"`
	Port    int      `json:"port"`
}

func LoadConfig(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err == nil {
		var cfg Config
		if err := json.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config: %w", err)
		}

		for i, s := range cfg.Symbols {
			cfg.Symbols[i] = strings.ToLower(strings.TrimSpace(s))
		}
		return &cfg, nil
	}

	symbolsStr := os.Getenv("CRYPTO_SYMBOLS")
	var symbols []string

	if symbolsStr != " " {
		symbols = strings.Split(symbolsStr, ",")
		for i, s := range symbols {
			symbols[i] = strings.ToLower(strings.TrimSpace(s))
		}
	} else {
		symbols = []string{"btcusdt", "ethusdt", "bnbusdt", "solusdt"}
	}
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		fmt.Sscanf(portStr, "%d", &port)
	}

	return &Config{
		Symbols: symbols,
		Port:    port,
	}, nil
}
