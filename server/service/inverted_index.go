package service

import (
	"context"
	"log"
	"math"
	"server/dao"
)

type IndexChoice uint8

const (
	word2article IndexChoice = iota
	word2author
	author2article
)

// Get all inverted-indexes of the words.
func getInvertedIndexes(fields []string, choice IndexChoice) (invertedIndexes []dao.InvertedIndex) {
	invertedIndexes = make([]dao.InvertedIndex, 0)
	for _, field := range fields {
		invertedIndex := getInvertedIndex(field, choice)
		invertedIndexes = append(invertedIndexes, invertedIndex)
	}
	return
}

// Try to get from cache(check by bloom filter) before get from mysql.
func getInvertedIndex(field string, choice IndexChoice) (invertedIndex dao.InvertedIndex) {
	if choice == word2article {
		if !dao.ArticleWordFilter.TestString(field) {
			return nil
		}

		var _invertedIndex string
		err := dao.WordToArticleCache.Get(context.Background(), field, &_invertedIndex)
		if err != nil {
			err = dao.DB.Model(&dao.WordToArticle{}).Where("word = ?", field).
				Select("indexes").Find(&_invertedIndex).Error
			if err != nil {
				return nil
			}
		}
		err = invertedIndex.UnSerialize(_invertedIndex)
		if err != nil {
			log.Println("service/inverted_index.go getInvertedIndex error:", err)
		}
		return
	} else if choice == word2author {
		if !dao.AuthorWordFilter.TestString(field) {
			return nil
		}

		var _invertedIndex string
		err := dao.WordToAuthorCache.Get(context.Background(), field, &_invertedIndex)
		if err != nil {
			err = dao.DB.Model(&dao.WordToAuthor{}).Where("word = ?", field).
				Select("indexes").Find(&_invertedIndex).Error
			if err != nil {
				return nil
			}
		}
		err = invertedIndex.UnSerialize(_invertedIndex)
		if err != nil {
			log.Println("service/inverted_index.go getInvertedIndex error:", err)
		}
		return
	} else if choice == author2article {
		var _invertedIndex string
		err := dao.AuthorToArticleRDB.Get(context.Background(), field).Scan(&_invertedIndex)
		if err == nil {
			err = invertedIndex.UnSerialize(_invertedIndex)
			if err != nil {
				log.Println("service/inverted_index.go getInvertedIndex error:", err)
			}
			return
		}
		temp := make([]uint64, 0)
		err = dao.DB.Model(&dao.AuthorToArticle{}).Where("id = ?", field).
			Select("article_id").Find(&temp).Error
		if err != nil {
			return nil
		}
		for _, articleId := range temp {
			invertedIndex.Add(articleId)
		}
		return
	} else {
		return nil
	}
}

func bm25(id uint64, invertedIndexes []dao.InvertedIndex, choice IndexChoice) (score float64) {
	var (
		totalCnt, wordCnt int64
		avgWordCnt        float64
		err               error
	)

	if choice == word2author {
		totalCnt = dao.AuthorCnt
		avgWordCnt = dao.AuthorAvgWordCnt
		wordCnt, err = getAuthorWordCount(id)
		if err != nil {
			log.Println("service/inverted_index.go bm25 error:", err)
			return
		}
	} else if choice == word2article {
		totalCnt = dao.ArticleCnt
		avgWordCnt = dao.ArticleAvgWordCnt
		wordCnt, err = getArticleWordCount(id)
		if err != nil {
			log.Println("service/inverted_index.go bm25 error:", err)
			return
		}
	} else {

		return 0
	}

	for _, invertedIndex := range invertedIndexes {
		var _score float64
		if _, ok := invertedIndex[id]; ok {
			_score = math.Log(1. + (float64(totalCnt)-float64(len(invertedIndex))+0.5)/
				(float64(len(invertedIndex))) + 0.5)
			// let k1 = 1.2, b = 0.75.
			_score *= float64(invertedIndex[id]) / float64(wordCnt) * 2.2
			_score /= float64(invertedIndex[id])/float64(wordCnt) +
				1.2*(1-0.75+0.75*float64(wordCnt)/avgWordCnt)
		}
		score += _score
	}

	return
}
