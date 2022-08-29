package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/jdkato/prose/tokenize"
	"github.com/jinzhu/inflection"
	jsoniter "github.com/json-iterator/go"
	"github.com/sajari/fuzzy"
	"log"
	"os"
	"regexp"
	"runtime"
	"server/dao"
	"strconv"
	"strings"
	"time"
)

var (
	json         = jsoniter.ConfigCompatibleWithStandardLibrary
	spellChecker *fuzzy.Model
	NumCpu       int
)

func Init() {
	NumCpu = runtime.NumCPU()

	spellChecker = fuzzy.NewModel()
	spellChecker.SetDepth(2)
	spellChecker.Train(SampleEnglish())
}

func textToWord(text string) (words []string) {
	if len(text) == 0 {
		return
	}

	_words := tokenize.TextToWords(text)
	for _, word := range _words {
		for len(word) > 0 && !('a' <= word[0] && word[0] <= 'z' ||
			'A' <= word[0] && word[0] <= 'Z' || '0' <= word[0] && word[0] <= '9') {
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
			word = strings.ToLower(word)
			t := spellChecker.SpellCheck(word)
			if len(t) != 0 {
				word = t
			}
			word = inflection.Singular(word)
			words = append(words, word)
		}
	}
	return
}

// Make sure the article id exist before use.
func getArticleAuthor(id uint64) (authorIDs []uint64) {
	sID := strconv.FormatUint(id, 10)

	_authorIDs := dao.ArticleToAuthorRDB.SMembers(context.Background(), sID).Val()
	if len(_authorIDs) == 0 {
		dao.DB.Model(&dao.ArticleToAuthor{}).Where("article_id = ?", id).
			Select("author_id").Find(&_authorIDs)

		go func() {
			err := dao.ArticleToAuthorRDB.SAdd(context.Background(), sID, _authorIDs).Err()
			if err != nil {
				log.Println("service/common.go getArticleAuthor error:", err)
			}
			dao.ArticleToAuthorRDB.Expire(context.Background(), sID, time.Minute)
		}()
	}

	authorIDs = make([]uint64, len(_authorIDs))
	for i := range _authorIDs {
		authorIDs[i], _ = strconv.ParseUint(_authorIDs[i], 10, 64)
	}

	return
}

// Make sure the article id exist before use.
func getAuthorArticle(id uint64) (articleIDs []uint64) {
	sID := strconv.FormatUint(id, 10)

	_authorIDs := dao.AuthorToArticleRDB.SMembers(context.Background(), sID).Val()
	if len(_authorIDs) == 0 {
		dao.DB.Model(&dao.AuthorToArticle{}).Where("author_id = ?", id).
			Select("article_id").Find(&_authorIDs)

		go func() {
			err := dao.AuthorToArticleRDB.SAdd(context.Background(), sID, _authorIDs).Err()
			if err != nil {
				log.Println("service/common.go getAuthorArticle error:", err)
			}
			dao.AuthorToArticleRDB.Expire(context.Background(), sID, time.Minute)
		}()
	}

	articleIDs = make([]uint64, len(_authorIDs))
	for i := range _authorIDs {
		articleIDs[i], _ = strconv.ParseUint(_authorIDs[i], 10, 64)
	}

	return
}

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
