package config

import (
	"github.com/mishankov/go-utlz/cliutils"
)

var BotSecret = cliutils.GetEnvOrDefault("BOT_SECRET", "secret")
var Port = cliutils.GetEnvOrDefault("PORT", ":4444")
