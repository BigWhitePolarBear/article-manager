package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

type WordToAuthor struct {
	Word    string `gorm:"primaryKey; type:varchar(100) not null"`
	Indexes string `gorm:"type:longtext"`
}

func (w *WordToAuthor) BeforeSave(tx *gorm.DB) (err error) {
	go w.deleteFromCache()
	return nil
}

func (w *WordToAuthor) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		w.deleteFromCache()
	}()
	return nil
}

func (w *WordToAuthor) BeforeUpdate(tx *gorm.DB) (err error) {
	go w.deleteFromCache()
	return nil
}

func (w *WordToAuthor) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		w.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (w *WordToAuthor) AfterFind(tx *gorm.DB) (err error) {
	go w.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (w *WordToAuthor) AfterCreate(tx *gorm.DB) (err error) {
	go w.saveIntoCache()

	return nil
}

func (w *WordToAuthor) saveIntoCache() {
	err := WordToAuthorCache.Set(&cache.Item{
		Key:   w.Word,
		Value: w.Indexes,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/word_author.go saveIntoCache error:", err)
	}
}

func (w *WordToAuthor) deleteFromCache() {
	retriedTimes := 0
retry:
	err := WordToAuthorCache.Delete(context.Background(), w.Word)
	if err != nil && retriedTimes < 5 {
		retriedTimes++
		goto retry
	}
}
