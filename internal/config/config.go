package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - структура конфига
type Config struct {
	Env          string         `yaml:"env" env-default:"local"`
	StoragePaths string         `yaml:"storage_paths" env-required:"true"`
	TokenTTL     time.Duration  `yaml:"token_ttl" env-required:"true"`
	Port         int            `yaml:"port" env-default:"8080"`
	SConfig      ServerConfig   `yaml:"server" env-required:"true"`
	DBConfig     DataBaseConfig `yaml:"db" env-required:"true"`
}

type ServerConfig struct {
	Port string `yaml:"port" env-default:"8080"`
	Host string `yaml:"host" env-default:"localhost"`
}

type DataBaseConfig struct {
	Port     string `yaml:"port" env-default:"5436"`
	Host     string `yaml:"host" env-default:"localhost"`
	Username string `yaml:"username" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-default:"5436"`
}

// MustLoad получает структуру конфига
func MustLoad() *Config {
	path := fetchConfigFlags()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config" + err.Error())
	}

	return &cfg
}

// fetchConfigFlags получает путь до конфига либо из флага командной строки либо через переменную окружения
func fetchConfigFlags() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
