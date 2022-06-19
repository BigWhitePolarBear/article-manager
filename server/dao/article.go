package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync/atomic"
	"time"
)

type Article struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Title     string         `gorm:"type:varchar(1500) not null"`
	Book      Book           `json:"-"`
	BookID    *uint64        `json:",omitempty"`
	Journal   Journal        `json:"-"`
	JournalID *uint64        `json:",omitempty"`
	Volume    string         `json:",omitempty" gorm:"type:varchar(50)"`
	Pages     string         `json:",omitempty" gorm:"type:varchar(50)"`
	Year      *uint16        `json:",omitempty"`
	EE        string         `json:",omitempty" gorm:"type:varchar(500)"`
}

func (a *Article) BeforeSave(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *Article) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *Article) BeforeUpdate(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *Article) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *Article) AfterFind(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *Article) AfterCreate(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	atomic.AddInt64(&authorCount, 1)

	return nil
}

func (a *Article) saveIntoCache() {
	jsonA, err := json.Marshal(*a)
	if err != nil {
		log.Println("dao/article.go saveIntoCache error:", err)
	}

	err = ArticleCache.Set(&cache.Item{
		Key:   strconv.FormatUint(a.ID, 10),
		Value: jsonA,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/article.go saveIntoCache error:", err)
	}
}

func (a *Article) deleteFromCache() {
	_ = ArticleCache.Delete(context.Background(), strconv.FormatUint(a.ID, 10))
}
