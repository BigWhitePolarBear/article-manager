package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type AuthorWordCount struct {
	ID    uint64 `gorm:"primaryKey"`
	Count uint64
}

func (a *AuthorWordCount) BeforeSave(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *AuthorWordCount) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *AuthorWordCount) BeforeUpdate(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *AuthorWordCount) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *AuthorWordCount) AfterFind(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *AuthorWordCount) AfterCreate(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

func (a *AuthorWordCount) saveIntoCache() {
	err := ArticleWordCntCache.Set(&cache.Item{
		Key:   strconv.FormatUint(a.ID, 10),
		Value: strconv.FormatUint(a.Count, 10),
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/author_word_count.go saveIntoCache error:", err)
	}
}

func (a *AuthorWordCount) deleteFromCache() {
	_ = AuthorCache.Delete(context.Background(), strconv.FormatUint(a.ID, 10))
}
