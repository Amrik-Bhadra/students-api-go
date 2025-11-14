package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addrress string  `yaml:"address" env-required:"true"`
}

// creating a struct
type Config struct {
	Env         string               `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string               `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"` // struct embedding
}

func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		// path not passed by env, try another way
		// check in arguments
		flags := flag.String("config", "", "path to the configurtion file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config is not set")
		}
	}


	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesnot exist: %s", configPath)
	}

	var config Config

	// serealize
	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatalf("Cannot Read Config File %s", err.Error())
	}

	return &config

}