package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string `yaml:"env"`
	HTTPServer struct {
		Address     string        `yaml:"address"`
		Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
		IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	} `yaml:"http_server"`
	Storage struct {
		MongoDB struct {
			DSN    string `yaml:"dsn"`
			DBName string `yaml:"db_name"`
		} `yaml:"mongodb"`
	} `yaml:"storage"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Load config path is failed")
	}

	//panic("Load config failed")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exists", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Parse config is failed: %s", err.Error())
	}
	return &cfg
}
