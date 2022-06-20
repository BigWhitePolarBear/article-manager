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

type Author struct {
	ID           uint64         `gorm:"primaryKey"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Name         string         `gorm:"type:varchar(100) not null; index:,length:10"`
	ArticleCount uint16         `gorm:"index:,sort:desc"`
}

func (a *Author) BeforeSave(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *Author) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *Author) BeforeUpdate(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *Author) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *Author) AfterFind(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *Author) AfterCreate(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	atomic.AddInt64(&AuthorCnt, 1)

	return nil
}

func (a *Author) saveIntoCache() {
	jsonB, err := json.Marshal(*a)
	if err != nil {
		log.Println("dao/author.go saveIntoCache error:", err)
	}

	err = AuthorCache.Set(&cache.Item{
		Key:   strconv.FormatUint(a.ID, 10),
		Value: jsonB,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/author.go saveIntoCache error:", err)
	}
}

func (a *Author) deleteFromCache() {
	_ = AuthorCache.Delete(context.Background(), strconv.FormatUint(a.ID, 10))
}
