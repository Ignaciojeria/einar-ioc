package configuration

import (
	"log/slog"
	"os"

	ioc "github.com/Ignaciojeria/einar-ioc"
	"github.com/joho/godotenv"
)

var _ = ioc.Registry(NewConf)

type Conf struct {
	Port      string
	ApiPrefix string
}

func NewConf() (Conf, error) {
	if err := godotenv.Load(); err != nil {
		slog.Warn(".env not found, loading environment from system.")
	}
	conf := Conf{
		Port:      os.Getenv("Port"),
		ApiPrefix: os.Getenv("ApiPrefix"),
	}
	if conf.ApiPrefix == "" {
		conf.ApiPrefix = "/api/"
	}
	if conf.Port == "" {
		conf.Port = "8080"
	}
	return conf, nil
}
