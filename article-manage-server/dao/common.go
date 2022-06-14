package dao

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title       string `gorm:"type:varchar(600) not null; index:,length:50"`
	Authors     string
	Journal     Journal `gorm:"foreignKey:JournalID"`
	JournalID   *uint
	Volume      string  `gorm:"type:varchar(50)"`
	Month       string  `gorm:"type:varchar(10)"`
	Year        *uint16 `gorm:"index"`
	CdRom       string  `gorm:"type:varchar(50)"`
	EE          string
	Publisher   Publisher `gorm:"foreignKey:PublisherID"`
	PublisherID *uint
	ISBN        string `gorm:"type:varchar(50)"`
}

type Author struct {
	gorm.Model
	Name         string `gorm:"type:varchar(100) not null; index:,length:10"`
	Articles     string
	ArticleCount uint16 `gorm:"index:,sort:desc"`
}

type Publisher struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
}

type Journal struct {
	gorm.Model
	Name string `gorm:"type:varchar(100)"`
}
