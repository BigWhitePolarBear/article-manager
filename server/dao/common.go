package dao

// All type in dao package got delay double deletion strategy for their cache.

import (
	"github.com/bits-and-blooms/bloom/v3"
	jsoniter "github.com/json-iterator/go"
)

type Variable struct {
	Key   string
	Value string
}

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	ArticleCnt        int64
	AuthorCnt         int64
	ArticleAvgWordCnt float32
	AuthorAvgWordCnt  float32

	ArticleIDFilter   *bloom.BloomFilter
	ArticleWordFilter *bloom.BloomFilter

	AuthorIDFilter   *bloom.BloomFilter
	AuthorWordFilter *bloom.BloomFilter
)

func Init() {
	DB.Model(&Article{}).Count(&ArticleCnt)
	DB.Model(&Author{}).Count(&AuthorCnt)

	err := DB.Model(&Variable{}).Where("`key` = ?", "ArticleAvgWordCnt").
		Select("value").Find(&ArticleAvgWordCnt).Error
	if err != nil {
		panic(err)
	}
	err = DB.Model(&Variable{}).Where("`key` = ?", "AuthorAvgWordCnt").
		Select("value").Find(&AuthorAvgWordCnt).Error
	if err != nil {
		panic(err)
	}

	// load the bloom filters
	var (
		jsonArticleIDFilter   string
		jsonArticleWordFilter string
		jsonAuthorIDFilter    string
		jsonAuthorWordFilter  string
	)

	err = DB.Model(&Variable{}).Where("`key` = ?", "ArticleIDFilter").
		Select("value").Find(&jsonArticleIDFilter).Error
	if err != nil {
		panic(err)
	}
	err = DB.Model(&Variable{}).Where("`key` = ?", "ArticleWordFilter").
		Select("value").Find(&jsonArticleWordFilter).Error
	if err != nil {
		panic(err)
	}
	err = DB.Model(&Variable{}).Where("`key` = ?", "AuthorIDFilter").
		Select("value").Find(&jsonAuthorIDFilter).Error
	if err != nil {
		panic(err)
	}
	err = DB.Model(&Variable{}).Where("`key` = ?", "AuthorWordFilter").
		Select("value").Find(&jsonAuthorWordFilter).Error
	if err != nil {
		panic(err)
	}

	ArticleIDFilter = &bloom.BloomFilter{}
	ArticleWordFilter = &bloom.BloomFilter{}

	AuthorIDFilter = &bloom.BloomFilter{}
	AuthorWordFilter = &bloom.BloomFilter{}

	err = ArticleIDFilter.UnmarshalJSON([]byte(jsonArticleIDFilter))
	if err != nil {
		panic(err)
	}
	err = ArticleWordFilter.UnmarshalJSON([]byte(jsonArticleWordFilter))
	if err != nil {
		panic(err)
	}
	err = AuthorIDFilter.UnmarshalJSON([]byte(jsonAuthorIDFilter))
	if err != nil {
		panic(err)
	}
	err = AuthorWordFilter.UnmarshalJSON([]byte(jsonAuthorWordFilter))
	if err != nil {
		panic(err)
	}
}
