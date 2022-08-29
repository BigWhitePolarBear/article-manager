package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func getInvertedIndexes(fields []string, choice IndexChoice) []dao.InvertedIndex {
	invertedIndexes := make([]dao.InvertedIndex, 0)
	for _, field := range fields {
		invertedIndex := getInvertedIndex(field, choice)
		invertedIndexes = append(invertedIndexes, invertedIndex)
	}
	return invertedIndexes
}

// Try to get from cache(check by bloom filter) before get from mysql, return nil if there is an error.
func getInvertedIndex(field string, choice IndexChoice) dao.InvertedIndex {
	invertedIndex := dao.InvertedIndex{}
	if choice == word2article {
		if !dao.ArticleWordFilter.TestString(field) {
			return nil
		}

		var _invertedIndex string
		err := dao.WordToArticleCache.Get(context.Background(), field, &_invertedIndex)
		if err != nil {
			wordToArticle := dao.WordToArticle{}
			// use struct to storage for gorm hook.
			err = dao.DB.Model(&dao.WordToArticle{}).Where("word = ?", field).
				Find(&wordToArticle).Error
			if err != nil {
				return nil
			}
			_invertedIndex = wordToArticle.Indexes
		}
		err = invertedIndex.UnSerialize(_invertedIndex)
		if err != nil {
			log.Println("service/inverted_index.go getInvertedIndex error:", err)
			return nil
		}
		return invertedIndex

	} else if choice == word2author {
		if !dao.AuthorWordFilter.TestString(field) {
			return nil
		}

		var _invertedIndex string
		err := dao.WordToAuthorCache.Get(context.Background(), field, &_invertedIndex)
		if err != nil {
			wordToAuthor := dao.WordToAuthor{}
			// use struct to storage for gorm hook.
			err = dao.DB.Model(&dao.WordToAuthor{}).Where("word = ?", field).
				Find(&wordToAuthor).Error
			if err != nil {
				return nil
			}
			_invertedIndex = wordToAuthor.Indexes
		}
		err = invertedIndex.UnSerialize(_invertedIndex)
		if err != nil {
			log.Println("service/inverted_index.go getInvertedIndex error:", err)
			return nil
		}
		return invertedIndex

	} else if choice == author2article {
		authorID, _ := strconv.ParseUint(field, 10, 64)
		articleIDs := getAuthorArticle(authorID)
		for _, articleId := range articleIDs {
			invertedIndex.Add(articleId)
		}
		return invertedIndex
	}

	// Wrong choice.
	return nil
}

// Get all inverted-indexes of the words directly from mysql with locks.
func getInvertedIndexesForUpdate(tx *gorm.DB, fields []string, choice IndexChoice) ([]dao.InvertedIndex, error) {
	invertedIndexes := make([]dao.InvertedIndex, 0)
	for _, field := range fields {
		invertedIndex, err := getInvertedIndexForUpdate(tx, field, choice)
		if err != nil {
			return nil, err
		}
		invertedIndexes = append(invertedIndexes, invertedIndex)
	}
	return invertedIndexes, nil
}

// Get inverted-index of the word directly from mysql with lock.
func getInvertedIndexForUpdate(tx *gorm.DB, field string, choice IndexChoice) (dao.InvertedIndex, error) {
	invertedIndex := dao.InvertedIndex{}
	if choice == word2article {
		if !dao.ArticleWordFilter.TestString(field) {
			return invertedIndex, nil
		}
		wordToArticle := dao.WordToArticle{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.WordToArticle{}).
			Where("word = ?", field).Find(&wordToArticle).Error
		if err != nil {
			return nil, err
		}
		err = invertedIndex.UnSerialize(wordToArticle.Indexes)
		if err != nil {
			return nil, err
		}
		return invertedIndex, nil

	} else if choice == word2author {
		if !dao.AuthorWordFilter.TestString(field) {
			return invertedIndex, nil
		}
		wordToAuthor := dao.WordToAuthor{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.WordToAuthor{}).
			Where("word = ?", field).Find(&wordToAuthor).Error
		if err != nil {
			return nil, err
		}
		err = invertedIndex.UnSerialize(wordToAuthor.Indexes)
		if err != nil {
			return nil, err
		}
		return invertedIndex, nil

	} else if choice == author2article {
		authorID, _ := strconv.ParseUint(field, 10, 64)
		articleIDs := getAuthorArticle(authorID)
		for _, articleId := range articleIDs {
			invertedIndex.Add(articleId)
		}
		return invertedIndex, nil
	}

	// Wrong choice.
	return nil, errors.New("wrong choice param in getInvertedIndexesForUpdate func")
}

// Save all inverted-indexes of the words.
// Words should be promised to exist since they were got form getInvertedIndexForUpdate func.
// Only support word2article and word2author.
func saveInvertedIndexes(tx *gorm.DB, words []string, invertedIndexes []dao.InvertedIndex, choice IndexChoice) error {
	for i := range words {
		err := saveInvertedIndex(tx, words[i], invertedIndexes[i], choice)
		if err != nil {
			return err
		}
	}
	return nil
}

// Save the inverted-index of a word.
// Word should be promised to exist since it was got form getInvertedIndexForUpdate func.
// Only support word2article and word2author.
func saveInvertedIndex(tx *gorm.DB, word string, invertedIndex dao.InvertedIndex, choice IndexChoice) error {
	if choice == word2article {
		if len(invertedIndex) == 0 {
			err := tx.Delete(&dao.WordToArticle{Word: word}).Error
			if err != nil {
				return err
			}
		} else {
			err := tx.Save(&dao.WordToArticle{Word: word, Indexes: invertedIndex.Serialize()}).Error
			if err != nil {
				return err
			}
		}
		return nil
	} else if choice == word2author {
		if len(invertedIndex) == 0 {
			err := tx.Delete(&dao.WordToAuthor{Word: word}).Error
			if err != nil {
				return err
			}
		} else {
			err := tx.Save(&dao.WordToAuthor{Word: word, Indexes: invertedIndex.Serialize()}).Error
			if err != nil {
				return err
			}
		}
		return nil
	}
	return errors.New("no such choice")
}

func bm25(id uint64, invertedIndexes []dao.InvertedIndex, choice IndexChoice) (score float32) {
	var (
		totalCnt   uint64
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
