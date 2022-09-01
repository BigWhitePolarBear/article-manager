package service

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"server/dao"
)

func Update(tmpArticle Article) (dao.Article, error) {
	id := tmpArticle.ID

	article := dao.Article{}

	if id == 0 {
		return article, errors.New("please input article id")
	}

	// Modify data in transaction.
	tx := dao.DB.Begin()

	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Find(&article).Error
	if err != nil {
		tx.Rollback()
		return article, err
	}

	article.Book.Name, article.Journal.Name, article.Volume = tmpArticle.Book, tmpArticle.Journal, tmpArticle.Volume
	article.Year, article.EE, article.Pages = &tmpArticle.Year, tmpArticle.EE, tmpArticle.Pages
	article.ID, article.Authors = tmpArticle.ID, tmpArticle.Authors

	// Lock global data.
	err = tx.Raw("select * from variables for update").Error
	if err != nil {
		tx.Rollback()
		return article, err
	}

	// Used when transaction need to be rolled back.
	oldArticleAvgWordCnt := dao.ArticleAvgWordCnt
	oldAuthorCnt, oldAuthorAvgWordCnt := dao.AuthorCnt, dao.AuthorAvgWordCnt

	if tmpArticle.Title != "" {
		oldTitleWords, newTitleWords := textToWord(article.Title), textToWord(tmpArticle.Title)

		// Update related inverted-indexes.
		oldTitleWordSet, newTitleWordSet := wordSlice2set(oldTitleWords), wordSlice2set(newTitleWords)
		interWordSet := wordSetIntersection(oldTitleWordSet, newTitleWordSet)
		delWordSet, addWordSet := wordSetDif(oldTitleWordSet, interWordSet), wordSetDif(newTitleWordSet, interWordSet)
		delWords, addWords := wordSet2Slice(delWordSet), wordSet2Slice(addWordSet)

		delWordIndexes, err := getInvertedIndexesForUpdate(tx, delWords, word2article)
		if err != nil {
			tx.Rollback()
			return article, err
		}
		for i := range delWordIndexes {
			delWordIndexes[i].Del(id)
		}
		err = saveInvertedIndexes(tx, delWords, delWordIndexes, word2article)
		if err != nil {
			tx.Rollback()
			return article, err
		}

		addWordIndexes, err := getInvertedIndexesForUpdate(tx, addWords, word2article)
		if err != nil {
			tx.Rollback()
			return article, err
		}
		for i := range addWordIndexes {
			addWordIndexes[i].Add(id)
		}
		err = saveInvertedIndexes(tx, addWords, addWordIndexes, word2article)
		if err != nil {
			tx.Rollback()
			return article, err
		}

		article.Title = tmpArticle.Title

		err = tx.Model(&dao.ArticleWordCount{ID: id}).Where("id = ?", id).
			Update("count", len(newTitleWords)).Error
		if err != nil {
			tx.Rollback()
			return article, err
		}

		// Update global data
		dao.ArticleAvgWordCnt = ((float32(dao.ArticleCnt) * dao.ArticleAvgWordCnt) -
			float32(len(oldTitleWords)+len(newTitleWords))) / float32(dao.ArticleCnt)
		err = tx.Model(&dao.Variable{}).Where("`key` = ?", "ArticleAvgWordCnt").
			Update("value", dao.ArticleAvgWordCnt).Error
		if err != nil {
			tx.Rollback()
			dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
			return article, err
		}
	}

	if len(tmpArticle.Authors) > 0 {
		oldAuthors := getArticleAuthor(id)
		newNames := tmpArticle.Authors
		// Try to lock these authors' data.
		for _, name := range newNames {
			tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.Author{}).Where("name = ?", name)
		}
		oldNames := make([]string, len(oldAuthors))
		for i, authorID := range oldAuthors {
			err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.Author{}).
				Where("id = ?", authorID).Select("name").Find(&oldNames[i]).Error
			if err != nil {
				tx.Rollback()
				dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
				return article, err
			}
		}
		oldNameSet, newNameSet := wordSlice2set(oldNames), wordSlice2set(newNames)
		interNameSet := wordSetIntersection(oldNameSet, newNameSet)
		delNameSet, addNameSet := wordSetDif(oldNameSet, interNameSet), wordSetDif(newNameSet, interNameSet)

		for name := range delNameSet {
			author := dao.Author{}
			err = tx.Model(&dao.Author{}).Where("name = ?", name).Find(&author).Error
			if err != nil {
				tx.Rollback()
				dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
				dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
				return article, err
			}

			author.ArticleCount--
			if author.ArticleCount == 0 {
				err = delAuthor(tx, author.ID)
				if err != nil {
					tx.Rollback()
					dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
					dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
					return article, err
				}
			} else {
				err = tx.Model(&dao.Author{ID: author.ID}).Where("id = ?", author.ID).
					Update("article_count", author.ArticleCount).Error
				if err != nil {
					tx.Rollback()
					dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
					dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
					return article, err
				}
			}

			// Delete connection.
			err = delConnection(tx, id, author.ID)
			if err != nil {
				tx.Rollback()
				dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
				dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
				return article, err
			}
		}

		for name := range addNameSet {
			author := dao.Author{}
			err = tx.Model(&dao.Author{}).Where("name = ?", name).Take(&author).Error
			if err != nil {
				// Check if it's a new author.
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = createAuthor(tx, name, id)
					if err != nil {
						tx.Rollback()
						dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
						dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
						return article, err
					}
				} else {
					tx.Rollback()
					dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
					dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
					return article, err
				}

			} else {
				author.ArticleCount++
				err = tx.Model(&dao.Author{ID: author.ID}).Where("id = ?", author.ID).
					Update("article_count", author.ArticleCount).Error
				if err != nil {
					tx.Rollback()
					dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
					dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
					return article, err
				}

				err = addConnection(tx, id, author.ID)
				if err != nil {
					tx.Rollback()
					dao.ArticleAvgWordCnt = oldArticleAvgWordCnt
					dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
					return article, err
				}
			}
		}
	}

	tx.Commit()

	return article, nil
}

type stringSet map[string]struct{}

func wordSlice2set(words []string) stringSet {
	set := stringSet{}
	for _, word := range words {
		set[word] = struct{}{}
	}
	return set
}

func wordSet2Slice(words stringSet) []string {
	slice := make([]string, 0, len(words))
	for word := range words {
		slice = append(slice, word)
	}
	return slice
}

func wordSetIntersection(A, B stringSet) stringSet {
	set := stringSet{}
	if len(A) > len(B) {
		A, B = B, A
	}
	for word := range A {
		if _, ok := B[word]; ok {
			set[word] = struct{}{}
		}
	}
	return set
}

// Return A-B
func wordSetDif(A, B stringSet) stringSet {
	set := stringSet{}
	for word := range A {
		if _, ok := B[word]; !ok {
			set[word] = struct{}{}
		}
	}
	return set
}
