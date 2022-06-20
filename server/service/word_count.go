package service

import (
	"context"
	"server/dao"
	"strconv"
)

func getArticleWordCount(id uint64) (count int64, err error) {
	err = dao.ArticleWordCntCache.Get(context.Background(), strconv.FormatUint(id, 10), &count)
	if err != nil {
		err = dao.DB.Model(&dao.ArticleWordCount{}).Where("id = ?", id).
			Select("count").Find(&count).Error
	}
	return
}

func getAuthorWordCount(id uint64) (count int64, err error) {
	err = dao.AuthorWordCntCache.Get(context.Background(), strconv.FormatUint(id, 10), &count)
	if err != nil {
		err = dao.DB.Model(&dao.AuthorWordCount{}).Where("id = ?", id).
			Select("count").Find(&count).Error
	}
	return
}
