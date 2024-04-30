package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App          App
		SqlDb        SqlDb
		CqlDb        CqlDb
		Cache        Cache
		AccessToken  AccessToken
		RefreshToken RefreshToken
	}

	App struct {
		Port             int
		RequestSizeLimit int
	}

	SqlDb struct {
		Host                string
		Port                int
		User                string
		Password            string
		DBName              string
		SSLMode             string
		TimeZone            string
		InitTableScriptPath string
	}

	CqlDb struct {
		Host                   string
		Port                   int
		KeySpace               string
		InitKeyspaceScriptPath string
		InitTableScriptPath    string
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
			Port:             viper.GetInt("app.server.port"),
			RequestSizeLimit: viper.GetInt("app.server.request_size_limit"),
		},

		SqlDb: SqlDb{
			Host:                viper.GetString("sql_database.host"),
			Port:                viper.GetInt("sql_database.port"),
			User:                viper.GetString("sql_database.user"),
			Password:            viper.GetString("sql_database.password"),
			DBName:              viper.GetString("sql_database.dbname"),
			SSLMode:             viper.GetString("sql_database.sslmode"),
			TimeZone:            viper.GetString("sql_database.timezone"),
			InitTableScriptPath: viper.GetString("sql_database.init_table_script_path"),
		},

		CqlDb: CqlDb{
			Host:                   viper.GetString("cql_database.host"),
			Port:                   viper.GetInt("cql_database.port"),
			KeySpace:               viper.GetString("cql_database.keyspace"),
			InitKeyspaceScriptPath: viper.GetString("cql_database.init_keyspace_script_path"),
			InitTableScriptPath:    viper.GetString("cql_database.init_table_script_path"),
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
