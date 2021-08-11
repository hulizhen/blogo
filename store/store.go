package store

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/hulizhen/blogo/config"
	"github.com/hulizhen/blogo/store/model"

	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	ArticleStore *model.ArticleStore
	AboutStore   *model.AboutStore
}

func New(cfg *config.Config) (*Store, error) {
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
		return nil, err
	}

	driver, err := mysql.WithInstance(db.DB, &mysql.Config{})
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithDatabaseInstance("file://store/migration", "mysql", driver)
	if err != nil {
		return nil, err
	}
	_ = m.Up()

	articleStore, err := model.NewArticleStore(db, cfg)
	if err != nil {
		return nil, err
	}

	aboutStore, err := model.NewAboutStore(cfg)
	if err != nil {
		return nil, err
	}

	return &Store{
		ArticleStore: articleStore,
		AboutStore:   aboutStore,
	}, nil
}
