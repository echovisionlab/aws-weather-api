package database

import (
	"fmt"
	"github.com/echovisionlab/aws-weather-api/pkg/constants"
	"os"
)

type Config struct {
	Name string
	User string
	Pass string
	Host string
	Port string
}

func (d *Config) ConnStr() string {
	return fmt.Sprintf("user=%v dbname=%v sslmode=disable password=%v host=%v port=%v",
		d.User,
		d.Name,
		d.Pass,
		d.Host,
		d.Port)
}

func WithEnvConfig() *Config {
	return &Config{
		Name: os.Getenv(constants.DatabaseName),
		User: os.Getenv(constants.DatabaseUser),
		Pass: os.Getenv(constants.DatabasePass),
		Host: os.Getenv(constants.DatabaseHost),
		Port: os.Getenv(constants.DatabasePort),
	}
}
