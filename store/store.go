package store

import (
	"blogo/config"
	"blogo/store/model"
	"fmt"

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
