package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"sync"
)

const (
	appPort    = "port"
	dbPort     = "db.port"
	dbUsername = "db.username"
	dbHost     = "db.host"
	dbName     = "db.name"
	dbPassword = "db.password"
)

type Repository struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
}

type Config struct {
	AppPort string `yaml:"port"`
	Repository
}

var instance *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		log.Println("initializing .yml file")
		if err := initConfig(); err != nil {
			panic(fmt.Errorf("panic while initializing .yml file, %v", err))
		}

		portUint, err := strconv.ParseUint(viper.GetString(dbPort), 10, 16)
		if err != nil {
			return
		}
		storageInstance := Repository{
			Host:     viper.GetString(dbHost),
			Port:     uint16(portUint),
			Username: viper.GetString(dbUsername),
			Password: viper.GetString(dbPassword),
			DBName:   viper.GetString(dbName),
		}
		instance = &Config{
			AppPort:    viper.GetString(appPort),
			Repository: storageInstance,
		}
	})

	return instance
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
