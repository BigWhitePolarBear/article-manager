package service

import (
	"context"
	"server/dao"
	"strconv"
)

func getArticleWordCount(id uint64) (count uint8, err error) {
	err = dao.ArticleWordCntCache.Get(context.Background(), strconv.FormatUint(id, 10), &count)
	if err != nil {
		wordCount := dao.ArticleWordCount{}
		err = dao.DB.Model(&dao.ArticleWordCount{}).Where("id = ?", id).Find(&wordCount).Error
		count = wordCount.Count
	}
	return
}

func getAuthorWordCount(id uint64) (count uint8, err error) {
	err = dao.AuthorWordCntCache.Get(context.Background(), strconv.FormatUint(id, 10), &count)
	if err != nil {
		wordCount := dao.AuthorWordCount{}
		err = dao.DB.Model(&dao.AuthorWordCount{}).Where("id = ?", id).Find(&wordCount).Error
		count = wordCount.Count
	}
	return
}
