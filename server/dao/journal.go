package dao

import (
	"context"
	"github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type Journal struct {
	ID        uint64         `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `gorm:"type:varchar(100); index:,length:10"`
}

func (j *Journal) BeforeSave(tx *gorm.DB) (err error) {
	j.deleteFromCache()
	return nil
}

func (j *Journal) AfterSave(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		j.deleteFromCache()
	}()
	return nil
}

func (j *Journal) BeforeUpdate(tx *gorm.DB) (err error) {
	j.deleteFromCache()
	return nil
}

func (j *Journal) AfterUpdate(tx *gorm.DB) (err error) {
	go func() {
		time.Sleep(200 * time.Millisecond)
		j.deleteFromCache()
	}()
	return nil
}

// AfterFind write into cache after search
func (j *Journal) AfterFind(tx *gorm.DB) (err error) {
	j.saveIntoCache()

	return nil
}

// AfterCreate write into cache after creation
func (j *Journal) AfterCreate(tx *gorm.DB) (err error) {
	j.saveIntoCache()

	return nil
}

func (j *Journal) saveIntoCache() {
	jsonJ, err := json.Marshal(*j)
	if err != nil {
		log.Println("dao/journal.go saveIntoCache error:", err)
	}

	err = JournalCache.Set(&cache.Item{
		Key:   strconv.FormatUint(j.ID, 10),
		Value: jsonJ,
		TTL:   time.Minute,
	})
	if err != nil {
		log.Println("dao/journal.go saveIntoCache error:", err)
	}
}

func (j *Journal) deleteFromCache() {
	_ = JournalCache.Delete(context.Background(), strconv.FormatUint(j.ID, 10))
}
