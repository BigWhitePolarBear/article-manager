package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

type WordToArticle struct {
	Word    string `gorm:"primaryKey; type:varchar(100) not null"`
	Indexes string `gorm:"type:longtext"`
}

func (w *WordToArticle) BeforeSave(tx *gorm.DB) (err error) {
	w.deleteFromCache()
	return nil
}

func (w *WordToArticle) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		w.deleteFromCache()
	}()
	return nil
}

func (w *WordToArticle) BeforeUpdate(tx *gorm.DB) (err error) {
	w.deleteFromCache()
	return nil
}

func (w *WordToArticle) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		w.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (w *WordToArticle) AfterFind(tx *gorm.DB) (err error) {
	w.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (w *WordToArticle) AfterCreate(tx *gorm.DB) (err error) {
	w.saveIntoCache()

	return nil
}

func (w *WordToArticle) saveIntoCache() {
	err := WordToArticleCache.Set(&cache.Item{
		Key:   w.Word,
		Value: w.Indexes,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/word_article.go saveIntoCache error:", err)
	}
}

func (w *WordToArticle) deleteFromCache() {
	_ = WordToArticleCache.Delete(context.Background(), w.Word)
}
