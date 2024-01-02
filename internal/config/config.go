package config

import (
	"flag"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	PgConn   string        `yaml:"pg_conn" env-required:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-default:"1h"`
	GRPC     GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path)
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		newPgConn, err := setHostname(cfg.PgConn, dbHost)
		if err != nil {
			panic("failed to set hostname for postgres connection: " + err.Error())
		}
		cfg.PgConn = newPgConn
	}

	return &cfg
}

func setHostname(addr, hostname string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	u.Host = hostname
	return u.String(), nil
}
