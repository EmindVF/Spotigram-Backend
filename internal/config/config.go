package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App          App
		Db           Db
		Cache        Cache
		AccessToken  AccessToken
		RefreshToken RefreshToken
	}

	App struct {
		Port int
	}

	Db struct {
		Host           string
		Port           int
		User           string
		Password       string
		DBName         string
		SSLMode        string
		TimeZone       string
		InitScriptPath string
	}

	Cache struct {
		RedisUrl string
	}

	AccessToken struct {
		PublicKey  string
		PrivateKey string
		ExpiresIn  time.Duration
		MaxAge     int
	}

	RefreshToken struct {
		PublicKey  string
		PrivateKey string
		ExpiresIn  time.Duration
		MaxAge     int
	}
)

// Returns the config object with fields
// parsed from the configuration file.
func GetConfig() Config {
	viper.SetConfigName("configuration")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %v", err))
	}

	return Config{
		App: App{
			Port: viper.GetInt("app.server.port"),
		},

		Db: Db{
			Host:           viper.GetString("database.host"),
			Port:           viper.GetInt("database.port"),
			User:           viper.GetString("database.user"),
			Password:       viper.GetString("database.password"),
			DBName:         viper.GetString("database.dbname"),
			SSLMode:        viper.GetString("database.sslmode"),
			TimeZone:       viper.GetString("database.timezone"),
			InitScriptPath: viper.GetString("database.initscriptpath"),
		},

		Cache: Cache{
			RedisUrl: viper.GetString("cache.redis_url"),
		},

		AccessToken: AccessToken{
			PrivateKey: viper.GetString("access_token.private_key"),
			PublicKey:  viper.GetString("access_token.public_key"),
			ExpiresIn:  viper.GetDuration("access_token.expires_in"),
			MaxAge:     viper.GetInt("access_token.max_age"),
		},

		RefreshToken: RefreshToken{
			PrivateKey: viper.GetString("refresh_token.private_key"),
			PublicKey:  viper.GetString("refresh_token.public_key"),
			ExpiresIn:  viper.GetDuration("refresh_token.expires_in"),
			MaxAge:     viper.GetInt("refresh_token.max_age"),
		},
	}
}
