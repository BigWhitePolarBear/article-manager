package service

import (
	"context"
	"log"
	"server/dao"
	"strconv"
)

func getArticle(id uint64, admin bool) (article dao.Article, ok bool) {
	// test with bloom filter
	if !dao.ArticleIDFilter.TestString(strconv.FormatUint(id, 10)) {
		return article, false
	}

	var jsonArticle []byte
	err := dao.ArticleCache.Get(context.Background(), strconv.FormatUint(id, 10), &jsonArticle)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonArticle, &article)
		if err != nil {
			log.Println("service/article.go getArticle error:", err)
			ok = false
		}
	} else {
		// cache missed
		err = dao.DB.Model(&dao.Article{}).Where("id = ?", id).Find(&article).Error
		if err != nil {
			log.Println("service/article.go getArticle error:", err)
			ok = false
		}
	}

	if !admin {
		article.ID = 0
	}
	return
}
