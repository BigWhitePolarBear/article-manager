package service

import (
	"errors"
	"github.com/jdkato/prose/tokenize"
	"server/dao"
	"sort"
	"strings"
)

type QueryType uint8

const (
	TitleQuery QueryType = iota
	YearQuery
	AuthorQuery
)

type IndexScore struct {
	ID    uint64
	Score float64
}

func SearchArticle(query map[QueryType]string, admin bool) (articles []dao.Article, err error) {
	invertedIndexes := make([]dao.InvertedIndex, 0)

	// get inverted-indexes by title queries.
	titleText, ok := query[TitleQuery]
	if ok {
		titleWords := queryTextToWord(titleText)

		invertedIndexes = getInvertedIndexes(titleWords, word2article)
		if 2*len(titleWords) < len(invertedIndexes) {
			return nil, errors.New("no match records for your title queries")
		}
	}

	// get inverted-indexes by authors.
	authorText, ok := query[AuthorQuery]
	if ok {
		authorWords := queryTextToWord(authorText)
		authorInvertedIndexes := getInvertedIndexes(authorWords, word2author)
		if err != nil {
			return nil, err
		}

		authorIndexes, err := dao.Union(authorInvertedIndexes)
		if err != nil {
			return nil, err
		}

		if len(authorIndexes) == 0 {
			return nil, errors.New("no match records for your author queries")
		}

		strAuthorIndexes := make([]string, len(authorIndexes))

		_invertedIndexes := getInvertedIndexes(strAuthorIndexes, author2article)
		if err != nil {
			return nil, err
		}

		invertedIndexes = append(invertedIndexes, _invertedIndexes...)
	}

	articleIndexes, err := dao.Union(invertedIndexes)
	if err != nil {
		return nil, err
	}
	indexScores := make([]IndexScore, len(articleIndexes))
	for i := range articleIndexes {
		indexScores[i].ID = articleIndexes[i]
		indexScores[i].Score = bm25(articleIndexes[i], invertedIndexes, word2article)
	}

	sort.Slice(indexScores, func(i, j int) bool {
		return indexScores[i].Score < indexScores[j].Score
	})

	for _, is := range indexScores {
		article, ok := getArticle(is.ID, admin)
		if ok {
			articles = append(articles, article)
		}
	}

	return
}

func GetAuthor(name string, limit, offset int) ([]Author, error) {
	searchSQL := authorsTable.Limit(limit).Offset(offset).Where("name = ?", name)
	_results := []dao.Author{}
	err := searchSQL.Find(&_results).Error

	if err != nil {
		return nil, err
	}

	results := make([]Author, 0, len(_results))

	// Turn work id to work title.
	for _, _result := range _results {
		worksBuilder := strings.Builder{}
		workIDs := strings.Fields(_result.Articles)
		for _, workID := range workIDs {
			var title string
			worksTable.Where("id = ?", workID).Select("title").Find(&title)
			if worksBuilder.Len() != 0 {
				worksBuilder.WriteString(", ")
			}
			worksBuilder.WriteString(title)
		}
		results = append(results, Author{_result.Name, worksBuilder.String(), _result.ArticleCount})
	}
	return results, nil
}

func GetTopAuthor(limit, offset int) (*[]Author, error) {
	results := make([]Author, limit)
	err := authorsTable.Select("name", "work_count").Order("work_count desc").Limit(limit).Offset(offset).Find(&results).Error
	if err != nil {
		return nil, err
	}
	return &results, nil
}

func queryTextToWord(text string) (words []string) {
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
			words = append(words, word)
		}
	}
	return
}
