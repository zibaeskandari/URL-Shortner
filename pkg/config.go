package pkg

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Token     TokenConfig     `mapstructure:"token"`
	Logging   LoggingConfig   `mapstructure:"logging"`
	ShortCode ShortCodeConfig `mapstructure:"short-code"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

type DatabaseConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	Database        string `mapstructure:"database"`
	SSLMode         string `mapstructure:"sslmode"`
	Dialect         string `mapstructure:"dialect"`
	MaxConns        int32  `mapstructure:"max_conns"`
	MinConns        int32  `mapstructure:"min_conns"`
	MaxConnLifetime int    `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime int    `mapstructure:"max_conn_idle_time"`
}

type TokenConfig struct{}

type LoggingConfig struct{}

type ShortCodeConfig struct {
	MinLength int    `mapstructure:"min-length"`
	MaxLength int    `mapstructure:"max-length"`
	Chars     string `mapstructure:"chars"`
	Digits    string `mapstructure:"digits"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	setDefaults(v) // we set a default so everything has a value

	v.SetEnvPrefix("APP")
	v.AutomaticEnv()                                   // this checks the current os environment and use if there is any of fields set there
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // since env names are SERVER_PORT this will map that to server.port in yml file

	setupEnvBindings(v) // viper does not handle dash in names properly

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("Error unmarshalling config, %s", err)
	}

	return &config, nil
}

func setupEnvBindings(v *viper.Viper) {
	// since it had dash in the name, viper failed to map it to env variable. Therefor, we manually bid it here
	_ = v.BindEnv("short-code.min-length", "APP_SHORT_CODE_MIN_LENGTH")
	_ = v.BindEnv("short-code.max-length", "APP_SHORT_CODE_MAX_LENGTH")
	_ = v.BindEnv("short-code.chars", "APP_SHORT_CODE_CHARS")
	_ = v.BindEnv("short-code.digits", "APP_SHORT_CODE_DIGITS")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 3000)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.env", "local")
	v.SetDefault("short-code.min-length", 5)
	v.SetDefault("short-code.max-length", 12)
	v.SetDefault("short-code.chars", "abdefghmnpqrstuvwxyz")
	v.SetDefault("short-code.digits", "23456789")
}
