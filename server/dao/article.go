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
	Title     string         `json:",omitempty" gorm:"type:varchar(1500) not null"`
	Authors   []string       `json:",omitempty" gorm:"-"`
	Book      Book
	BookID    *uint64
	Journal   Journal
	JournalID *uint64
	Volume    string  `json:",omitempty" gorm:"type:varchar(50)"`
	Pages     string  `json:",omitempty" gorm:"type:varchar(50)"`
	Year      *uint16 `json:",omitempty"`
	EE        string  `json:",omitempty" gorm:"type:varchar(500)"`
}

func (a *Article) BeforeSave(tx *gorm.DB) error {
	go a.deleteFromCache()
	return nil
}

func (a *Article) AfterSave(tx *gorm.DB) error {
	go func() {
		time.Sleep(500 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *Article) BeforeUpdate(tx *gorm.DB) error {
	go a.deleteFromCache()
	return nil
}

func (a *Article) AfterUpdate(tx *gorm.DB) error {
	go func() {
		time.Sleep(500 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *Article) BeforeDelete(tx *gorm.DB) error {
	go a.deleteFromCache()
	return nil
}

func (a *Article) AfterDelete(tx *gorm.DB) error {
	go func() {
		time.Sleep(500 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *Article) AfterFind(tx *gorm.DB) error {
	go a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *Article) AfterCreate(tx *gorm.DB) error {
	go a.saveIntoCache()

	atomic.AddUint64(&AuthorCnt, 1)

	return nil
}

func (a *Article) saveIntoCache() {
	sID := strconv.FormatUint(a.ID, 10)

	article := *a

	// Retrieved all from mysql.
	if len(a.Authors) > 0 {
		article.Title = ""
		jsonA, err := json.Marshal(article)
		if err != nil {
			log.Println("dao/article.go saveIntoCache error:", err)
		}

		err = ArticleCache.Set(&cache.Item{
			Key:   sID,
			Value: jsonA,
			TTL:   time.Minute,
		})
		if err != nil {
			log.Println("dao/article.go saveIntoCache error:", err)
		}
	}

	// Only retrieved title.
	err := TitleCache.Set(&cache.Item{
		Key:   sID,
		Value: a.Title,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/article.go saveIntoCache error:", err)
	}
}

func (a *Article) deleteFromCache() {
	sID := strconv.FormatUint(a.ID, 10)

	retriedTimes1, retriedTimes2 := 0, 0
retry1:
	err := ArticleCache.Delete(context.Background(), sID)
	if err != nil && retriedTimes1 < 5 {
		retriedTimes1++
		goto retry1
	}
retry2:
	err = TitleCache.Delete(context.Background(), sID)
	if err != nil && retriedTimes2 < 5 {
		retriedTimes2++
		goto retry2
	}
}
