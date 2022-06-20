package main

import _ "gorm.io/gorm"

type WordToArticle struct {
	Word    string `gorm:"primaryKey; type:varchar(100) not null"`
	Indexes string `gorm:"type:longtext"`
}

type WordToAuthor struct {
	Word    string `gorm:"primaryKey; type:varchar(100) not null"`
	Indexes string `gorm:"type:longtext"`
}

type ArticleWordCount struct {
	ID    uint64 `gorm:"primaryKey"`
	Count uint64
}

type AuthorWordCount struct {
	ID    uint64 `gorm:"primaryKey"`
	Count uint64
}

func init() {
	err := DB.AutoMigrate(&WordToAuthor{}, &WordToArticle{}, &ArticleWordCount{}, &AuthorWordCount{})
	if err != nil {
		panic(err)
	}
}
