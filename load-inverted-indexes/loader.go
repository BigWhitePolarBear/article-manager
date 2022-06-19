package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bits-and-blooms/bloom/v3"
	"github.com/jdkato/prose/tokenize"
	"github.com/jinzhu/inflection"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func SampleEnglish() []string {
	var out []string
	file, err := os.Open("/project/article-manager/data/fuzzy/big.txt")
	if err != nil {
		fmt.Println(err)
		return out
	}
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	// Count the words.
	count := 0
	for scanner.Scan() {
		exp, _ := regexp.Compile("[a-zA-Z]+")
		words := exp.FindAll([]byte(scanner.Text()), -1)
		for _, word := range words {
			if len(word) > 1 {
				out = append(out, strings.ToLower(string(word)))
				count++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}

	return out
}

func articleLoader() {
	defer wg.Done()

	articleIDFilter := bloom.NewWithEstimates(1e7, 0.1)
	articleWordFilter := bloom.NewWithEstimates(1e7, 0.1)

	var i uint64 = 1
	for ; i <= 2423332; i++ {
		receiver := struct {
			Id    uint64
			Title string
		}{}
		err := DB.Table("articles").Where("id = ?", i).Select("id", "title").Find(&receiver).Error
		if err != nil {
			log.Println(err)
			continue
		}
		articleIDFilter.AddString(strconv.FormatUint(receiver.Id, 16))

		words := tokenize.TextToWords(receiver.Title)
		var wordCnt uint64
		for _, word := range words {
			for len(word) > 0 && (word[0] == '\'' || word[0] == '-' || word[0] == '/' || word[0] == '.' ||
				word[0] == '`' || word[0] == '*' || word[0] == '+' || word[0] == '=' || word[0] == '^' ||
				word[0] == '\\' || word[0] == ',' || word[0] == '_' || word[0] == '|' || word[0] == '~') {
				word = word[1:]
			}
			if len(word) == 0 || len(word) == 1 {
				continue
			} else if len(word) == 2 && (word == "of" || word == "to" || word == "it" || word == "as" ||
				word == "or" || word == "in" || word == "on" || word == "'s" || word == "``") {
				continue
			} else if len(word) == 3 && (word == "and" || word == "for" || word == "its" || word == "the") {
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
				wordCnt++
				word = strings.ToLower(word)
				t := spellChecker.SpellCheck(word)
				if len(t) != 0 {
					word = t
				}
				word = inflection.Singular(word)
				articleWordFilter.AddString(word)
				WordToArticleRDB.HIncrBy(context.Background(), word, strconv.FormatUint(receiver.Id, 16), 1)
			}
		}
		err = DB.Model(&ArticleWordCount{}).Create(&ArticleWordCount{ID: i, Count: wordCnt}).Error
		if err != nil {
			log.Println(err)
		}
	}

	// Get all words and their indexes
	keys, err := WordToArticleRDB.Keys(context.Background(), "*").Result()
	if err != nil {
		panic(err)
	}

	for _, word := range keys {
		index := InvertedIndex{}
		indexes, err := WordToArticleRDB.HKeys(context.Background(), word).Result()
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
			_cnt, err := WordToArticleRDB.HGet(context.Background(), word, indexes[i]).Result()
			if err != nil {
				log.Println(err)
				continue
			}
			cnt, err := strconv.Atoi(_cnt)
			if err != nil {
				log.Println(err)
				continue
			}
			index[id] = cnt
		}

		DB.Model(&WordToArticle{}).Create(&WordToArticle{Word: word, Indexes: index.Serialize()})

		WordToArticleRDB.Del(context.Background(), word)
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

func authorLoader() {
	defer wg.Done()

	authorIDFilter := bloom.NewWithEstimates(1e7, 0.1)
	authorWordFilter := bloom.NewWithEstimates(1e7, 0.1)

	var i uint64 = 1
	for ; i <= 1540683; i++ {
		receiver := struct {
			Id   uint64
			Name string
		}{}
		err := DB.Table("authors").Where("id = ?", i).Select("id", "name").Find(&receiver).Error
		if err != nil {
			log.Println(err)
			continue
		}
		authorIDFilter.AddString(strconv.FormatUint(receiver.Id, 16))

		words := tokenize.TextToWords(receiver.Name)
		var wordCnt uint64
		for _, word := range words {
			for len(word) > 0 && (word[0] == '\'' || word[0] == '-' || word[0] == '/' || word[0] == '.' ||
				word[0] == '`' || word[0] == '*' || word[0] == '+' || word[0] == '=' || word[0] == '^' ||
				word[0] == '\\' || word[0] == ',' || word[0] == '_' || word[0] == '|' || word[0] == '~') {
				word = word[1:]
			}
			if len(word) == 0 || len(word) == 1 {
				continue
			} else {
				wordCnt++
				word = strings.ToLower(word)
				t := spellChecker.SpellCheck(word)
				if len(t) != 0 {
					word = t
				}
				authorWordFilter.AddString(word)
				WordToAuthorRDB.HIncrBy(context.Background(), word, strconv.FormatUint(receiver.Id, 16), 1)
			}
		}

		err = DB.Model(&AuthorWordCount{}).Create(&AuthorWordCount{ID: i, Count: wordCnt}).Error
		if err != nil {
			log.Println(err)
		}
	}

	// Get all words and their indexes
	keys, err := WordToAuthorRDB.Keys(context.Background(), "*").Result()
	if err != nil {
		panic(err)
	}

	for _, word := range keys {
		index := InvertedIndex{}
		indexes, err := WordToAuthorRDB.HKeys(context.Background(), word).Result()
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
			_cnt, err := WordToAuthorRDB.HGet(context.Background(), word, indexes[i]).Result()
			if err != nil {
				log.Println(err)
				continue
			}
			cnt, err := strconv.Atoi(_cnt)
			if err != nil {
				log.Println(err)
				continue
			}
			index[id] = cnt
		}

		DB.Model(&WordToAuthor{}).Create(&WordToAuthor{Word: word, Indexes: index.Serialize()})

		WordToAuthorRDB.Del(context.Background(), word)
	}

	// Storage bloom filter into database
	filterJson, err := authorIDFilter.MarshalJSON()
	if err != nil {
		panic(err)
	}
	DB.Table("variables").Create(variable{"AuthorIDFilter", string(filterJson)})

	filterJson, err = authorWordFilter.MarshalJSON()
	if err != nil {
		panic(err)
	}
	DB.Table("variables").Create(variable{"AuthorWordFilter", string(filterJson)})
}
