package main

import (
	"strings"
)

func record() {
	defer wg.Done()

	for {
		select {
		case article := <-articleMatch2Record:
			_authors := <-authorsMatch2Record
			authors := strings.Split(_authors, ", ")
			DB.Model(Article{}).Create(&article)

			for _, author := range authors {
				if author == "" {
					continue
				}
				tempAuthor := Author{}
				DB.Model(&Author{}).Order("name").Where("name = ?", author).Find(&tempAuthor)
				if tempAuthor.ID == 0 {
					tempAuthor.Name = author
					tempAuthor.ArticleCount = 1
					DB.Model(&Author{}).Create(&tempAuthor)
				} else {
					tempAuthor.ArticleCount++
					DB.Model(&Author{}).Where("id = ?", tempAuthor.ID).Update("article_count", tempAuthor.ArticleCount)
				}
				DB.Model(&ArticleToAuthor{}).Create(ArticleToAuthor{ArticleID: article.ID, AuthorID: tempAuthor.ID})
				DB.Model(&AuthorToArticle{}).Create(AuthorToArticle{AuthorID: tempAuthor.ID, ArticleID: article.ID})
			}

		case <-matchOK:
			return
		}
	}
}
