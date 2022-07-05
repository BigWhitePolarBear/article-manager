package service

import (
	"context"
	"fmt"
	"log"
	"server/dao"
	"sort"
	"strconv"
	"strings"
	"time"
)

func getArticle(id uint64, admin bool) (article dao.Article, ok bool) {
	var jsonArticle []byte
	err := dao.ArticleCache.Get(context.Background(), strconv.FormatUint(id, 10), &jsonArticle)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonArticle, &article)
		if err != nil {
			log.Println("service/article.go getArticle error:", err)
			ok = false
		}
		article.Title = getTitle(id)
	} else {
		// cache missed
		err = dao.DB.Model(&dao.Article{}).Where("id = ?", id).Find(&article).Error
		if err != nil {
			log.Println("service/article.go getArticle error:", err)
			ok = false
		}
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
	err := dao.TitleCache.Get(context.Background(), strconv.FormatUint(id, 10), &title)
	if err != nil {
		article := dao.Article{}
		err = dao.DB.Model(&dao.Article{ID: id}).Where("id = ?", id).
			Select("id", "title").Find(&article).Error
		if err != nil {
			log.Println("service/article.go getTitle error:", err)
		}
	}

	return
}

func getBook(id uint64) (book dao.Book) {
	var jsonBook []byte
	err := dao.BookCache.Get(context.Background(), strconv.FormatUint(id, 10), &jsonBook)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonBook, &book)
		if err != nil {
			log.Println("service/article.go getBook error:", err)
		}
	} else {
		// cache missed
		err = dao.DB.Model(&dao.Book{}).Where("id = ?", id).Find(&book).Error
		if err != nil {
			log.Println("service/article.go getBook error:", err)
		}
	}

	return
}

func getJournal(id uint64) (journal dao.Journal) {
	var jsonJournal []byte
	err := dao.JournalCache.Get(context.Background(), strconv.FormatUint(id, 10), &jsonJournal)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonJournal, &journal)
		if err != nil {
			log.Println("service/article.go getJournal error:", err)
		}
	} else {
		// cache missed
		err = dao.DB.Model(&dao.Journal{}).Where("id = ?", id).Find(&journal).Error
		if err != nil {
			log.Println("service/article.go getJournal error:", err)
		}
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

		err := dao.ArticleResRDB.Expire(context.Background(), key, time.Minute).Err()
		if err != nil {
			log.Println("service/article.go cacheArticleRes error:", err)
		}
	}

	return
}
