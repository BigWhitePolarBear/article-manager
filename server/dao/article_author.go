package dao

import (
	"context"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type ArticleToAuthor struct {
	ArticleID uint64 `gorm:"primaryKey;autoIncrement:false"`
	AuthorID  uint64 `gorm:"primaryKey;autoIncrement:false"`
}

func (a *ArticleToAuthor) BeforeSave(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *ArticleToAuthor) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *ArticleToAuthor) BeforeUpdate(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *ArticleToAuthor) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *ArticleToAuthor) AfterFind(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *ArticleToAuthor) AfterCreate(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

func (a *ArticleToAuthor) saveIntoCache() {
	err := ArticleToAuthorRDB.SAdd(context.Background(),
		strconv.FormatUint(a.ArticleID, 10), strconv.FormatUint(a.AuthorID, 10)).Err()
	if err != nil {
		log.Println("dao/article_author.go saveIntoCache error:", err)
	}

	err = ArticleToAuthorRDB.Expire(context.Background(), strconv.FormatUint(a.ArticleID, 10), time.Minute).Err()
	if err != nil {
		log.Println("dao/article_author.go saveIntoCache error:", err)
	}
}

func (a *ArticleToAuthor) deleteFromCache() {
	_ = ArticleToAuthorRDB.Del(context.Background(), strconv.FormatUint(a.ArticleID, 10))
}
