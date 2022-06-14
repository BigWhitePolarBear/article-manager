package main

import (
	"gorm.io/gorm"
)

type Article struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `gorm:"type:varchar(600) not null"`
	Authors   string
	Journal   Journal
	JournalID *uint64
	Volume    string `gorm:"type:varchar(50)"`
	Year      *uint16
	EE        string
}

type Author struct {
	ID           uint64         `gorm:"primaryKey"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"type:varchar(100) not null; index:,length:10"`
	Articles     string
	ArticleCount uint16 `gorm:"index:,sort:desc"`
}

type Journal struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `gorm:"type:varchar(100)"`
}
