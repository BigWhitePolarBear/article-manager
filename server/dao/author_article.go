package dao

import (
	"context"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type AuthorToArticle struct {
	AuthorID  uint64 `gorm:"primaryKey;autoIncrement:false"`
	ArticleID uint64 `gorm:"primaryKey;autoIncrement:false"`
}

func (a *AuthorToArticle) BeforeSave(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *AuthorToArticle) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

func (a *AuthorToArticle) BeforeUpdate(tx *gorm.DB) (err error) {
	a.deleteFromCache()
	return nil
}

func (a *AuthorToArticle) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		a.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (a *AuthorToArticle) AfterFind(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (a *AuthorToArticle) AfterCreate(tx *gorm.DB) (err error) {
	a.saveIntoCache()

	return nil
}

func (a *AuthorToArticle) saveIntoCache() {
	err := AuthorToArticleRDB.SAdd(context.Background(),
		strconv.FormatUint(a.AuthorID, 10), strconv.FormatUint(a.ArticleID, 10)).Err()
	if err != nil {
		log.Println("dao/author_article.go saveIntoCache error:", err)
	}
}

func (a *AuthorToArticle) deleteFromCache() {
	_ = AuthorToArticleRDB.Del(context.Background(), strconv.FormatUint(a.AuthorID, 10))
}
