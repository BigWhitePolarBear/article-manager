package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type ArticleWordCount struct {
	ID    uint64 `gorm:"primaryKey"`
	Count uint64
}

func (a *ArticleWordCount) BeforeSave(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *ArticleWordCount) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *ArticleWordCount) BeforeUpdate(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *ArticleWordCount) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *ArticleWordCount) AfterFind(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *ArticleWordCount) AfterCreate(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

func (a *ArticleWordCount) saveIntoCache() {
	err := ArticleWordCntCache.Set(&cache.Item{
		Key:   strconv.FormatUint(a.ID, 10),
		Value: strconv.FormatUint(a.Count, 10),
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/article_word_count.go saveIntoCache error:", err)
	}
}

func (a *ArticleWordCount) deleteFromCache() {
	_ = AuthorCache.Delete(context.Background(), strconv.FormatUint(a.ID, 10))
}
