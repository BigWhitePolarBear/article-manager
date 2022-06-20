package main

import (
	"gorm.io/gorm"
)

type Article struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `gorm:"type:varchar(1500) not null"`
	Book      Book
	BookID    *uint64
	Journal   Journal
	JournalID *uint64
	Volume    string `gorm:"type:varchar(50)"`
	Pages     string `gorm:"type:varchar(50)"`
	Year      *uint16
	EE        string `gorm:"type:varchar(500)"`
}

type Author struct {
	ID           uint64         `gorm:"primaryKey"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"type:varchar(100) not null; index:,length:10"`
	ArticleCount uint16         `gorm:"index:,sort:desc"`
}

type Journal struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"type:varchar(100); index:,length:10"`
}

type Book struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"type:varchar(100); index:,length:10"`
}

type ArticleToAuthor struct {
	ArticleID uint64 `gorm:"primaryKey;autoIncrement:false"`
	AuthorID  uint64 `gorm:"primaryKey;autoIncrement:false"`
}

type AuthorToArticle struct {
	AuthorID  uint64 `gorm:"primaryKey;autoIncrement:false"`
	ArticleID uint64 `gorm:"primaryKey;autoIncrement:false"`
}
