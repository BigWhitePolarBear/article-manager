package main

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	WordToArticleRDB *redis.Client
	WordToAuthorRDB  *redis.Client
)

const (
	numCommonRDB = iota
	numWordToArticleRDB
	numWordToAuthorRDB
)

func init() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:zxc05020519@tcp(localhost:3306)/"+
		"article_manager?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local"),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&WordToArticle{}, &WordToAuthor{})
	if err != nil {
		panic(err)
	}

	WordToArticleRDB = redis.NewClient(&redis.Options{
		Addr:     ":7000",
		DB:       numWordToArticleRDB,
		Password: "zxc05020519",
	})

	WordToAuthorRDB = redis.NewClient(&redis.Options{
		Addr:     ":7000",
		DB:       numWordToAuthorRDB,
		Password: "zxc05020519",
	})
}
