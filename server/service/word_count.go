package service

import (
	"context"
	"golang.org/x/sync/singleflight"
	"log"
	"server/dao"
	"strconv"
	"time"
)

var (
	articleWordGroup singleflight.Group
	authorWordGroup  singleflight.Group
)

func getArticleWordCount(id uint64) (count uint8, err error) {
	sID := strconv.FormatUint(id, 10)

	err = dao.ArticleWordCntCache.Get(context.Background(), sID, &count)
	if err != nil {
		// cache missed
		_count, _err, _ := articleWordGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				articleWordGroup.Forget(sID)
			}()

			wordCount := dao.ArticleWordCount{}
			err = dao.DB.Model(&dao.ArticleWordCount{}).Where("id = ?", id).Find(&wordCount).Error
			return wordCount.Count, err
		})
		err = _err
		if err != nil {
			log.Println("service/word_count.go getArticleWordCount error:", err)
			return
		}

		count = _count.(uint8)
	}
	return
}

func getAuthorWordCount(id uint64) (count uint8, err error) {
	sID := strconv.FormatUint(id, 10)

	err = dao.AuthorWordCntCache.Get(context.Background(), sID, &count)
	if err != nil {
		// cache missed
		_count, _err, _ := authorWordGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				authorWordGroup.Forget(sID)
			}()

			wordCount := dao.AuthorWordCount{}
			err = dao.DB.Model(&dao.AuthorWordCount{}).Where("id = ?", id).Find(&wordCount).Error
			return wordCount.Count, err
		})
		err = _err
		if err != nil {
			log.Println("service/word_count.go getAuthorWordCount error:", err)
			return
		}

		count = _count.(uint8)
	}
	return
}
