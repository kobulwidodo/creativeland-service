package config

import (
	"go-clean/src/lib/midtrans"
	"go-clean/src/lib/sql"
)

type Application struct {
	Meta     ApplicationMeta
	SQL      sql.Config
	Midtrans midtrans.Config
}

type ApplicationMeta struct {
	Title       string `mapstructure:"META_TITLE"`
	Description string `mapstructure:"META_DESCRIPTION"`
	Host        string `mapstructure:"META_HOST"`
	BasePath    string `mapstructure:"META_BASEPATH"`
	Version     string `mapstructure:"META_VERSION"`
}

func Init() Application {
	return Application{}
}
