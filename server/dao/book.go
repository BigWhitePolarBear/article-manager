package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type Book struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:",omitempty" gorm:"type:varchar(100); index:,length:10"`
}

func (b *Book) BeforeSave(tx *gorm.DB) (err error) {
	go b.deleteFromCache()
	return nil
}

func (b *Book) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		b.deleteFromCache()
	}()
	return nil
}

func (b *Book) BeforeUpdate(tx *gorm.DB) (err error) {
	go b.deleteFromCache()
	return nil
}

func (b *Book) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		b.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (b *Book) AfterFind(tx *gorm.DB) (err error) {
	go b.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (b *Book) AfterCreate(tx *gorm.DB) (err error) {
	go b.saveIntoCache()

	return nil
}

func (b *Book) saveIntoCache() {
	jsonB, err := json.Marshal(*b)
	if err != nil {
		log.Println("dao/book.go saveIntoCache error:", err)
	}

	err = BookCache.Set(&cache.Item{
		Key:   strconv.FormatUint(b.ID, 10),
		Value: jsonB,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/book.go saveIntoCache error:", err)
	}
}

func (b *Book) deleteFromCache() {
	retriedTimes := 0
retry:
	err := BookCache.Delete(context.Background(), strconv.FormatUint(b.ID, 10))
	if err != nil && retriedTimes < 5 {
		retriedTimes++
		goto retry
	}
}
