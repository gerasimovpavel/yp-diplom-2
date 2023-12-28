package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"strconv"
	"strings"
	"time"
)

const (
	appPrefix                                   = "BLUFF_"
	defaultHTTPPort               string        = "8080"
	defaultHTTPRWTimeout          time.Duration = 30 * time.Second
	defaultHTTPMaxHeaderMegabytes int           = 1
	defaultJWTSigningKey          string        = "123456"
	defaultAccessTokenTTL         time.Duration = 15 * time.Minute
	defaultRefreshTokenTTL        time.Duration = 24 * time.Hour * 30
)

type (
	Config struct {
		configPath string
		Database   MongoDB
		HTTP       HTTPConfig
		Auth       AuthConfig
		CacheTTL   time.Duration
	}

	MongoDB struct {
		URL      string `yaml:"url" env:"url"`
		Port     string `yaml:"port" env:"port"`
		User     string `yaml:"user" env:"user"`
		Password string `yaml:"password" env:"password"`
		Db       string `yaml:"database" env:"db"`
	}

	SQLite struct {
		URI      string `yaml:"uri" env:"uri"`
		User     string `yaml:"user" env:"user"`
		Password string `yaml:"password" env:"password"`
		Db       string `yaml:"database" env:"db"`
	}

	AuthConfig struct {
		JWT                    JWTConfig
		PasswordSalt           string `yaml:"password_salt" env:"salt"`
		VerificationCodeLength int    `yaml:"verification_code_length" env:"codelength"`
	}

	JWTConfig struct {
		AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env:"accessttl"`
		RefreshTokenTTL time.Duration `yaml:"refresh-token-ttl" env:"refreshttl"`
		SigningKey      string        `yaml:"signing_key" env:"signkey"`
	}

	HTTPConfig struct {
		Host               string        `yaml:"host" env:"host"`
		Port               string        `yaml:"port" env:"port"`
		ReadTimeout        time.Duration `yaml:"read_timeout" env:"readtimeout"`
		WriteTimeout       time.Duration `yaml:"write_timeout" env:"writetimeout"`
		MaxHeaderMegabytes int           `yaml:"max_header_megabytes" env:"maxheadermb"`
	}
)

func Init(configPath string) (*Config, error) {
	var cfg Config
	cfg.configPath = configPath
	cfg.defaults()
	if err := cfg.parseConfig(); err != nil {
		return nil, err
	}
	if err := cfg.parseEnv(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) defaults() {
	cfg.HTTP.Host = "0.0.0.0"
	cfg.HTTP.Port = defaultHTTPPort
	cfg.HTTP.WriteTimeout = defaultHTTPRWTimeout
	cfg.HTTP.ReadTimeout = defaultHTTPRWTimeout
	cfg.HTTP.MaxHeaderMegabytes = defaultHTTPMaxHeaderMegabytes
	cfg.Auth.JWT.AccessTokenTTL = defaultAccessTokenTTL
	cfg.Auth.JWT.RefreshTokenTTL = defaultRefreshTokenTTL
	cfg.Auth.JWT.SigningKey = defaultJWTSigningKey
}

func (cfg *Config) parseConfig() error {
	var k = koanf.New(".")
	if err := k.Load(file.Provider("configs/settings.yml"), yaml.Parser()); err != nil {
		return err
	}

	err := k.UnmarshalWithConf("", cfg, koanf.UnmarshalConf{Tag: "yaml"})
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) parseEnv() error {
	var k = koanf.New(".")

	if err := k.Load(env.ProviderWithValue(appPrefix, ".", func(s string, v string) (string, interface{}) {

		key := strings.Replace(strings.ToLower(strings.TrimPrefix(s, appPrefix)), "_", ".", -1)
		if value, err := strconv.Atoi(v); err == nil {
			return key, value
		}
		return key, v
	}), nil); err != nil {
		return err
	}
	err := k.UnmarshalWithConf("", cfg, koanf.UnmarshalConf{Tag: "env"})
	if err != nil {
		return err
	}
	return nil
}
