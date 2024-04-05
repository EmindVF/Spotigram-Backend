package config

import (
	"spotigram/internal/config"
)

var Cfg *config.Config

// Sets package's config to the privided config.
func SetupConfig(c *config.Config) {
	Cfg = c
}
