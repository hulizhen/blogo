package main

import (
	"fmt"
	"log"

	"blogo/config"
	"blogo/router"
	"blogo/service/observer"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg := config.New(config.ConfigFilePath)

	dsn := fmt.Sprintf(
		"%v:%v@tcp(%v:%d)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		cfg.Mysql.Username,
		cfg.Mysql.Password,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Database,
	)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Panicf("Failed to open database with error: %v.", err)
	}

	o, err := observer.NewRepoObserver(db, cfg.Website.BlogRepoPath)
	if err != nil {
		log.Panicf("Failed to create repo observer with error: %v.", err)
	}
	o.Run()

	r := router.New(cfg, db)
	r.Run()
}
