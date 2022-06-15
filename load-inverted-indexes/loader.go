package main

import (
	"context"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/jdkato/prose/tokenize"
	"log"
	"strconv"
	"strings"
)

func articleLoader() {
	defer wg.Done()

	articleIDFilter := bloom.NewWithEstimates(1e7, 0.01)
	articleWordFilter := bloom.NewWithEstimates(1e7, 0.01)

	rows, err := DB.Table("articles").Select("id", "title").Rows()
	if err != nil {
		panic(err)
	}

	var (
		id    uint64
		title string
	)
	for rows.Next() {
		err = rows.Scan(&id, &title)
		if err != nil {
			log.Println(err)
			continue
		}
		articleIDFilter.AddString(strconv.FormatUint(id, 16))
		words := tokenize.TextToWords(title)
		for _, word := range words {
			if len(word) == 1 {
				continue
			} else if len(word) == 2 && (word == "of" || word == "to" || word == "it" ||
				word == "as" || word == "or" || word == "in") {
				continue
			} else if len(word) == 3 && (word == "and" || word == "for" || word == "its") {
				continue
			} else if len(word) == 4 && (word == "with" || word == "when" || word == "that" ||
				word == "this") {
				continue
			} else if len(word) == 5 && (word == "while" || word == "about" || word == "their" ||
				word == "those" || word == "these") {
				continue
			} else if len(word) == 6 && (word == "across" || word == "inside") {
				continue
			} else {
				word = strings.ToLower(word)
				articleWordFilter.AddString(word)
				WordToArticleRDB.SAdd(context.Background(), word, id)
			}
		}
	}

	// Get all words and their indexes
	keys, err := WordToArticleRDB.Keys(context.Background(), "*").Result()
	if err != nil {
		panic(err)
	}

	for _, word := range keys {
		s := set{}
		indexes, err := WordToArticleRDB.SMembers(context.Background(), word).Result()
		if err != nil {
			log.Println(err)
			continue
		}

		for i := range indexes {
			id, err := strconv.ParseUint(indexes[i], 16, 64)
			if err != nil {
				log.Println(err)
				continue
			}
			s.put(id)
		}

		DB.Model(&WordToArticle{}).Create(&WordToArticle{Word: word, Indexes: s.serialize()})
	}

	// Storage bloom filter into database
	filterJson, err := articleIDFilter.MarshalJSON()
	if err != nil {
		panic(err)
	}
	DB.Table("variables").Create(variable{"ArticleIDFilter", string(filterJson)})

	filterJson, err = articleWordFilter.MarshalJSON()
	if err != nil {
		panic(err)
	}
	DB.Table("variables").Create(variable{"ArticleWordFilter", string(filterJson)})
}
