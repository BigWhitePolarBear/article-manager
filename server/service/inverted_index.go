package service

import (
	"context"
	"log"
	"math"
	"server/dao"
	"strconv"
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
	invertedIndex = dao.InvertedIndex{}
	if choice == word2article {
		if !dao.ArticleWordFilter.TestString(field) {
			return
		}

		var _invertedIndex string
		err := dao.WordToArticleCache.Get(context.Background(), field, &_invertedIndex)
		if err != nil {
			wordToArticle := dao.WordToArticle{}
			// use struct to storage for gorm hook.
			err = dao.DB.Model(&dao.WordToArticle{}).Where("word = ?", field).
				Find(&wordToArticle).Error
			if err != nil {
				return
			}
			_invertedIndex = wordToArticle.Indexes
		}
		err = invertedIndex.UnSerialize(_invertedIndex)
		if err != nil {
			log.Println("service/inverted_index.go getInvertedIndex error:", err)
		}
		return
	} else if choice == word2author {
		if !dao.AuthorWordFilter.TestString(field) {
			return
		}

		var _invertedIndex string
		err := dao.WordToAuthorCache.Get(context.Background(), field, &_invertedIndex)
		if err != nil {
			wordToAuthor := dao.WordToAuthor{}
			// use struct to storage for gorm hook.
			err = dao.DB.Model(&dao.WordToAuthor{}).Where("word = ?", field).
				Find(&wordToAuthor).Error
			if err != nil {
				return
			}
			_invertedIndex = wordToAuthor.Indexes
		}
		err = invertedIndex.UnSerialize(_invertedIndex)
		if err != nil {
			log.Println("service/inverted_index.go getInvertedIndex error:", err)
		}
		return
	} else if choice == author2article {
		authorID, _ := strconv.ParseUint(field, 10, 64)
		articleIDs := getAuthorArticle(authorID)
		for _, articleId := range articleIDs {
			invertedIndex.Add(articleId)
		}
		return
	} else {
		return
	}
}

func bm25(id uint64, invertedIndexes []dao.InvertedIndex, choice IndexChoice) (score float32) {
	var (
		totalCnt   int64
		wordCnt    uint8
		avgWordCnt float32
		err        error
	)

	if choice == word2article {
		totalCnt = dao.ArticleCnt
		avgWordCnt = dao.ArticleAvgWordCnt
		wordCnt, err = getArticleWordCount(id)
		if err != nil {
			log.Println("service/inverted_index.go bm25 error:", err)
			return
		}
	} else if choice == word2author {
		totalCnt = dao.AuthorCnt
		avgWordCnt = dao.AuthorAvgWordCnt
		wordCnt, err = getAuthorWordCount(id)
		if err != nil {
			log.Println("service/inverted_index.go bm25 error:", err)
			return
		}
	} else {

		return 0
	}

	for _, invertedIndex := range invertedIndexes {
		var _score float32
		floatLen := float32(len(invertedIndex))
		_score = float32(math.Log(float64(1. + (float32(totalCnt)-floatLen+0.5)/(floatLen+0.5))))
		// let k1 = 1.2, b = 0.75.
		floatCnt := float32(invertedIndex[id])
		floatWordCnt := float32(wordCnt)
		_score *= floatCnt / floatWordCnt * 2.2
		_score /= floatCnt/floatWordCnt + 0.3 + 0.9*floatWordCnt/avgWordCnt

		score += _score
	}

	return
}
