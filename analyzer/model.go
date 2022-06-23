package main

type Variable struct {
	Key   string
	Value string
}

type ArticleWordCount struct {
	ID    uint64 `gorm:"primaryKey"`
	Count uint64
}

type AuthorWordCount struct {
	ID    uint64 `gorm:"primaryKey"`
	Count uint64
}

type WordToArticle struct {
	Word    string `gorm:"primaryKey; type:varchar(100) not null"`
	Indexes string `gorm:"type:longtext"`
}

type WordToAuthor struct {
	Word    string `gorm:"primaryKey; type:varchar(100) not null"`
	Indexes string `gorm:"type:longtext"`
}
