package config

import (
	"spotigram/internal/config"
)

var Cfg *config.Config

func SetupConfig(c *config.Config) {
	Cfg = c
}
