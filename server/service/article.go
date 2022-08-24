package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"server/dao"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	articleGroup singleflight.Group
	titleGroup   singleflight.Group
	bookGroup    singleflight.Group
	journalGroup singleflight.Group
)

func getArticle(id uint64, admin bool) (article dao.Article, ok bool) {
	sID := strconv.FormatUint(id, 10)

	var jsonArticle []byte
	err := dao.ArticleCache.Get(context.Background(), sID, &jsonArticle)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonArticle, &article)
		if err != nil {
			log.Println("service/article.go getArticle error:", err)
			return article, false
		}
		article.Title = getTitle(id)
	} else {
		// cache missed
		_article, _err, _ := articleGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				articleGroup.Forget(sID)
			}()

			_article := dao.Article{}
			err = dao.DB.Model(&dao.Article{}).Where("id = ?", id).Find(&_article).Error
			return _article, err
		})
		err = _err
		if err != nil {
			log.Println("service/article.go getArticle error:", err)
			return article, false
		}

		article = _article.(dao.Article)
	}

	authorIDs := getArticleAuthor(id)
	article.Authors = make([]string, len(authorIDs))
	for i := range authorIDs {
		article.Authors[i] = getName(authorIDs[i])
	}

	if article.BookID != nil {
		article.Book = getBook(*article.BookID)
		if !admin {
			article.Book.ID = 0
		}
	}
	if article.JournalID != nil {
		article.Journal = getJournal(*article.JournalID)
		if !admin {
			article.Journal.ID = 0
		}
	}

	if !admin {
		article.ID = 0
	}

	return article, true
}

func getTitle(id uint64) (title string) {
	sID := strconv.FormatUint(id, 10)

	err := dao.TitleCache.Get(context.Background(), sID, &title)
	if err != nil {
		_title, _err, _ := titleGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				articleGroup.Forget(sID)
			}()

			_article := dao.Article{}
			err = dao.DB.Model(&dao.Article{ID: id}).Where("id = ?", id).
				Select("id", "title").Find(&_article).Error
			return _article.Title, err
		})
		err = _err
		if err != nil {
			log.Println("service/article.go getTitle error:", err)
			return
		}

		title = _title.(string)
	}

	return
}

func getBook(id uint64) (book dao.Book) {
	sID := strconv.FormatUint(id, 10)

	var jsonBook []byte
	err := dao.BookCache.Get(context.Background(), sID, &jsonBook)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonBook, &book)
		if err != nil {
			log.Println("service/article.go getBook error:", err)
			return
		}
	} else {
		// cache missed
		_book, _err, _ := bookGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				bookGroup.Forget(sID)
			}()

			_book := dao.Book{}
			err = dao.DB.Model(&dao.Book{ID: id}).Where("id = ?", id).
				Select("id", "name").Find(&_book).Error
			return _book, err
		})
		err = _err
		if err != nil {
			log.Println("service/article.go getBook error:", err)
			return
		}

		book = _book.(dao.Book)
	}

	return
}

func getJournal(id uint64) (journal dao.Journal) {
	sID := strconv.FormatUint(id, 10)

	var jsonJournal []byte
	err := dao.JournalCache.Get(context.Background(), sID, &jsonJournal)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonJournal, &journal)
		if err != nil {
			log.Println("service/article.go getJournal error:", err)
			return
		}
	} else {
		// cache missed
		_journal, _err, _ := journalGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				journalGroup.Forget(sID)
			}()

			_journal := dao.Journal{}
			err = dao.DB.Model(&dao.Journal{ID: id}).Where("id = ?", id).
				Select("id", "name").Find(&_journal).Error
			return _journal, err
		})
		err = _err
		if err != nil {
			log.Println("service/article.go getJournal error:", err)
			return
		}

		journal = _journal.(dao.Journal)
	}

	return
}

func getCachedArticleRes(titleWords, authorWords, notWords []string, page int, admin bool) (articles []dao.Article, ok bool) {
	builder := strings.Builder{}
	sort.Strings(titleWords)
	sort.Strings(authorWords)
	sort.Strings(notWords)
	for _, word := range titleWords {
		builder.WriteString(word)
	}
	for _, word := range authorWords {
		builder.WriteString(word)
	}
	for _, word := range notWords {
		builder.WriteString(word)
	}
	basicKey := builder.String()

	sIDs, err := dao.ArticleResRDB.LRange(context.Background(),
		fmt.Sprintf("%s:%d", basicKey, page), 0, -1).Result()
	if err != nil {
		return nil, false
	}

	exist := dao.ArticleResRDB.Exists(context.Background(), basicKey+":1").Val()
	if exist == 0 {
		return nil, false
	}

	for _, sID := range sIDs {
		id, _ := strconv.ParseUint(sID, 10, 64)
		article, ok := getArticle(id, admin)
		if ok {
			articles = append(articles, article)
		}
	}

	return articles, true
}

func cacheArticleRes(titleWords, authorWords, notWords []string, indexScores []IndexScore) {
	builder := strings.Builder{}
	sort.Strings(titleWords)
	sort.Strings(authorWords)
	sort.Strings(notWords)
	for _, word := range titleWords {
		builder.WriteString(word)
	}
	for _, word := range authorWords {
		builder.WriteString(word)
	}
	for _, word := range notWords {
		builder.WriteString(word)
	}
	basicKey := builder.String()

	for page := 1; page <= (len(indexScores)+49)/50; page++ {
		key := fmt.Sprintf("%s:%d", basicKey, page)
		for i := (page - 1) * 50; i < page*50 && i < len(indexScores); i++ {
			// store 10 base int for redis
			err := dao.ArticleResRDB.RPush(context.Background(), key,
				strconv.FormatUint(indexScores[i].ID, 10)).Err()
			if err != nil {
				log.Println("service/article.go cacheArticleRes error:", err)
			}
		}

		// Retry 5 times.
		var expireErr error
		expireErr = errors.New("")
		for i := 0; expireErr != nil && i < 5; i++ {
			expireErr = dao.ArticleResRDB.Expire(context.Background(), key, time.Minute).Err()
		}
		if expireErr != nil {
			log.Println("service/article.go cacheArticleRes error:", expireErr)
		}
	}

	return
}
