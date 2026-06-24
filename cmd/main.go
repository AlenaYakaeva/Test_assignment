package main

import (
	"exampleApp/internal"
	"exampleApp/internal/repository/db"
	"exampleApp/internal/server"
	"exampleApp/internal/service/wallet"
	"fmt"
)

func main() {
	cfg := internal.ReadConfig()
	repo, err := db.New(cfg.DBDSN)
	if err != nil {
		panic(err) //TODO обработать ошибку
	}
	if err = db.RunMigrations(cfg.DBDSN); err != nil {
		panic(err)
	}

	walService := wallet.New(repo)

	srv := server.New(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), walService)

	if err := srv.Run(); err != nil {
		panic(err)
	}
}
