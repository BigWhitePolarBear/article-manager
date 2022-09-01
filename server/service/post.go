package service

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"server/dao"
)

type Article struct {
	ID      uint64
	Title   string
	Authors []string `form:"Author"`
	Book    string
	Journal string
	Volume  string
	Pages   string
	EE      string
	Year    uint16
}

func Create(tmpArticle Article) (dao.Article, error) {
	newArticle := dao.Article{}
	newArticle.Title, newArticle.Book.Name = tmpArticle.Title, tmpArticle.Book
	newArticle.Journal.Name, newArticle.Volume = tmpArticle.Journal, tmpArticle.Volume
	newArticle.Pages, newArticle.EE, newArticle.Year = tmpArticle.Pages, tmpArticle.EE, &tmpArticle.Year

	// Modify data in transaction.
	tx := dao.DB.Begin()

	// Lock global data.
	err := tx.Raw("select * from variables for update").Error
	if err != nil {
		tx.Rollback()
		return newArticle, err
	}

	err = tx.Model(&dao.Article{}).Create(&newArticle).Error
	if err != nil {
		tx.Rollback()
		return newArticle, err
	}

	// Update related authors' data or create new authors.
	oldAuthorCnt, oldAuthorAvgWordCnt := dao.AuthorCnt, dao.AuthorAvgWordCnt
	for _, name := range tmpArticle.Authors {
		author := dao.Author{}
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.Author{}).
			Where("name = ?", name).Take(&author).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = createAuthor(tx, name, newArticle.ID)
				if err != nil {
					tx.Rollback()
					dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
					return newArticle, err
				}
			} else {
				tx.Rollback()
				dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
				return newArticle, err
			}

		} else {
			err = tx.Model(&dao.Author{ID: author.ID}).Where("id = ?", author.ID).
				Update("article_count", author.ArticleCount+1).Error
			if err != nil {
				tx.Rollback()
				dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
				return newArticle, err
			}

			// Create connection.
			err = addConnection(tx, newArticle.ID, author.ID)
			if err != nil {
				tx.Rollback()
				dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
				return newArticle, err
			}
		}
	}

	// Add this article id in related inverted indexes.
	titleWords := textToWord(newArticle.Title)

	err = tx.Create(&dao.ArticleWordCount{ID: newArticle.ID, Count: uint8(len(titleWords))}).Error
	if err != nil {
		tx.Rollback()
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return newArticle, err
	}

	indexes, err := getInvertedIndexesForUpdate(tx, titleWords, word2article)
	if err != nil {
		tx.Rollback()
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return newArticle, err
	}
	for i := range indexes {
		indexes[i].Add(newArticle.ID)
	}
	err = saveInvertedIndexes(tx, titleWords, indexes, word2article)
	if err != nil {
		tx.Rollback()
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return newArticle, err
	}

	// Update global data.
	oldArticleCnt, oldArticleAvgWordCnt := dao.ArticleCnt, dao.ArticleAvgWordCnt
	dao.ArticleCnt++
	dao.ArticleAvgWordCnt = ((float32(dao.ArticleCnt) * dao.ArticleAvgWordCnt) - float32(len(titleWords))) /
		float32(dao.ArticleCnt)

	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "ArticleCnt").
		Update("value", dao.ArticleCnt).Error
	if err != nil {
		tx.Rollback()
		dao.ArticleCnt, dao.ArticleAvgWordCnt = oldArticleCnt, oldArticleAvgWordCnt
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return newArticle, err
	}
	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "ArticleAvgWordCnt").
		Update("value", dao.ArticleAvgWordCnt).Error
	if err != nil {
		tx.Rollback()
		dao.ArticleCnt, dao.ArticleAvgWordCnt = oldArticleCnt, oldArticleAvgWordCnt
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return newArticle, err
	}

	tx.Commit()
	return newArticle, nil
}

// Created connection in this func.
func createAuthor(tx *gorm.DB, name string, articleID uint64) error {
	author := dao.Author{}
	// Create a new author.
	author.Name = name
	author.ArticleCount = 1
	err := tx.Model(&dao.Author{}).Create(&author).Error
	if err != nil {
		return err
	}

	// Create connection.
	err = addConnection(tx, articleID, author.ID)
	if err != nil {
		return err
	}

	// Add this author id in related inverted indexes.
	nameWords := textToWord(name)

	err = tx.Create(&dao.AuthorWordCount{ID: author.ID, Count: uint8(len(nameWords))}).Error
	if err != nil {
		return err
	}

	indexes, err := getInvertedIndexesForUpdate(tx, nameWords, word2author)
	if err != nil {
		return err
	}
	for i := range indexes {
		indexes[i].Add(author.ID)
	}
	err = saveInvertedIndexes(tx, nameWords, indexes, word2author)
	if err != nil {
		return err
	}

	// Update global data.
	dao.AuthorCnt++
	dao.AuthorAvgWordCnt = ((float32(dao.AuthorCnt) * dao.AuthorAvgWordCnt) - float32(len(nameWords))) /
		float32(dao.AuthorCnt)

	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "AuthorAvgWordCnt").
		Update("value", dao.AuthorAvgWordCnt).Error
	if err != nil {
		return err
	}
	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "AuthorCnt").
		Update("value", dao.AuthorCnt).Error
	if err != nil {
		return err
	}

	return nil
}

func addConnection(tx *gorm.DB, articleID, authorID uint64) error {
	err := tx.Create(&dao.ArticleToAuthor{ArticleID: articleID, AuthorID: authorID}).Error
	if err != nil {
		return err
	}
	err = tx.Create(&dao.AuthorToArticle{AuthorID: authorID, ArticleID: articleID}).Error
	if err != nil {
		return err
	}

	return nil
}
