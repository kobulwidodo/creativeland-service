package main

import (
	"fmt"
	"go-clean/src/business/domain"
	"go-clean/src/business/usecase"
	"go-clean/src/handler/rest"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/configreader"
	"go-clean/src/lib/midtrans"
	"go-clean/src/lib/sql"
	"go-clean/src/utils/config"

	_ "go-clean/docs/swagger"
)

// @contact.name   Rakhmad Giffari Nurfadhilah
// @contact.url    https://fadhilmail.tech/
// @contact.email  rakhmadgiffari14@gmail.com

// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization

const (
	configFile string = "./etc/cfg/app.env"
)

func main() {
	cfg := config.Init()
	configReader := configreader.Init(configreader.Options{
		ConfigFile: configFile,
	})

	// init meta config
	configReader.ReadConfig(&cfg.Meta)

	// init sql config
	configReader.ReadConfig(&cfg.SQL)

	// init midtrans config
	configReader.ReadConfig(&cfg.Midtrans)

	fmt.Printf("%#v", cfg)

	auth := auth.Init()

	midtrans := midtrans.Init(cfg.Midtrans)

	db := sql.Init(cfg.SQL)

	d := domain.Init(db, midtrans)

	uc := usecase.Init(auth, d)

	r := rest.Init(cfg.Meta, configReader, uc, auth)

	r.Run()
}
