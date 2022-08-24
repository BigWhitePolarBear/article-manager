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
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Name         string         `json:",omitempty" gorm:"type:varchar(100) not null; index:,length:10"`
	Articles     []string       `gorm:"-"`
	ArticleCount uint16         `gorm:"index:,sort:desc"`
}

func (a *Author) BeforeSave(tx *gorm.DB) (err error) {
	go a.deleteFromCache()
	return nil
}

func (a *Author) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *Author) BeforeUpdate(tx *gorm.DB) (err error) {
	go a.deleteFromCache()
	return nil
}

func (a *Author) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(500 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *Author) AfterFind(tx *gorm.DB) (err error) {
	go a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *Author) AfterCreate(tx *gorm.DB) (err error) {
	go a.saveIntoCache()

	atomic.AddInt64(&AuthorCnt, 1)

	return nil
}

func (a *Author) saveIntoCache() {
	sID := strconv.FormatUint(a.ID, 10)

	// Retrieved all from mysql.
	if len(a.Articles) > 0 {
		author := *a
		author.Name = ""
		jsonA, err := json.Marshal(author)
		if err != nil {
			log.Println("dao/author.go saveIntoCache error:", err)
		}

		err = AuthorCache.Set(&cache.Item{
			Key:   sID,
			Value: jsonA,
			TTL:   time.Minute,
		})
		if err != nil {
			log.Println("dao/author.go saveIntoCache error:", err)
		}
	}

	// Only retrieved name from mysql.
	if a.Name != "" {
		err := NameCache.Set(&cache.Item{
			Key:   sID,
			Value: a.Name,
			TTL:   time.Minute,
		})
		if err != nil {
			log.Println("dao/author.go saveIntoCache error:", err)
		}
	}
}

func (a *Author) deleteFromCache() {
	sID := strconv.FormatUint(a.ID, 10)

	retriedTimes1, retriedTimes2 := 0, 0
retry1:
	err := AuthorCache.Delete(context.Background(), sID)
	if err != nil && retriedTimes1 < 5 {
		retriedTimes1++
		goto retry1
	}
retry2:
	err = NameCache.Delete(context.Background(), sID)
	if err != nil && retriedTimes2 < 5 {
		retriedTimes2++
		goto retry2
	}
}
