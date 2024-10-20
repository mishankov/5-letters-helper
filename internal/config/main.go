package config

import (
	"os"
	"strings"
)

func envOrDefault(key, def string) string {
	value := os.Getenv(key)

	if len(value) == 0 {
		return def
	} else {
		return strings.TrimSpace(value)
	}
}

var BotSecret = envOrDefault("BOT_SECRET", "secret")
var Port = envOrDefault("PORT", ":4444")
