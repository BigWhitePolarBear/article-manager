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

type variable struct {
	Key   string
	Value string
}