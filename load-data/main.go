package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	read2match          = make(chan []byte, 16)
	articleMatch2Record = make(chan Article, 16)
	authorsMatch2Record = make(chan string, 16)
	readOK              = make(chan struct{})
	matchOK             = make(chan struct{})

	wg sync.WaitGroup
	DB *gorm.DB
)

func main() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:zxc05020519@tcp(localhost:3306)/"+
		"article_search_server?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local"),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	if err != nil {
		panic(err)
	}

	err = DB.AutoMigrate(&Article{}, &Author{}, &Journal{},
		&Book{}, &ArticleToAuthor{}, &AuthorToArticle{})
	if err != nil {
		panic(err)
	}
	wg.Add(3)

	go read()
	go match()
	go record()

	wg.Wait()
}
