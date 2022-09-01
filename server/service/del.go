package service

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"server/dao"
)

func DelArticle(id uint64) (dao.Article, error) {
	oldArticle := dao.Article{}

	// Modify data in transaction.
	tx := dao.DB.Begin()

	// Lock global data.
	err := tx.Raw("select * from variables for update").Error
	if err != nil {
		tx.Rollback()
		return oldArticle, err
	}

	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.Article{}).
		Where("id = ?", id).Find(&oldArticle).Error
	if err != nil {
		tx.Rollback()
		return oldArticle, err
	}

	err = tx.Delete(&dao.Article{ID: id}).Error
	if err != nil {
		tx.Rollback()
		return oldArticle, err
	}

	// Delete information on related authors.
	oldAuthorCnt, oldAuthorAvgWordCnt := dao.AuthorCnt, dao.AuthorAvgWordCnt
	authorIDs := getArticleAuthor(id)
	var tmpCnt uint16
	for _, authorID := range authorIDs {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Model(&dao.Author{}).
			Where("id = ?", authorID).Select("article_count").Find(&tmpCnt).Error
		if err != nil {
			tx.Rollback()
			dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
			return oldArticle, err
		}
		if tmpCnt != 1 {
			err = tx.Model(&dao.Author{ID: authorID}).Where("id = ?", authorID).
				Update("article_count", tmpCnt-1).Error
			if err != nil {
				tx.Rollback()
				dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
				return oldArticle, err
			}
		} else {
			// This author got no pub now, del it.
			err = delAuthor(tx, authorID)
			if err != nil {
				tx.Rollback()
				dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
				return oldArticle, err
			}
		}

		// Delete connection.
		err = delConnection(tx, id, authorID)
		if err != nil {
			tx.Rollback()
			dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
			return oldArticle, err
		}
	}

	// Del this article id in related inverted indexes.
	titleWords := textToWord(oldArticle.Title)

	err = tx.Delete(&dao.ArticleWordCount{ID: oldArticle.ID}).Error
	if err != nil {
		tx.Rollback()
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return oldArticle, err
	}

	indexes, err := getInvertedIndexesForUpdate(tx, titleWords, word2article)
	if err != nil {
		tx.Rollback()
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return oldArticle, err
	}
	for i := range indexes {
		indexes[i].Del(id)
	}
	err = saveInvertedIndexes(tx, titleWords, indexes, word2article)
	if err != nil {
		tx.Rollback()
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return oldArticle, err
	}

	// Update global data.
	oldArticleCnt, oldArticleAvgWordCnt := dao.ArticleCnt, dao.ArticleAvgWordCnt
	dao.ArticleCnt--
	dao.ArticleAvgWordCnt = ((float32(dao.ArticleCnt) * dao.ArticleAvgWordCnt) - float32(len(titleWords))) /
		float32(dao.ArticleCnt)

	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "ArticleCnt").
		Update("value", dao.ArticleCnt).Error
	if err != nil {
		tx.Rollback()
		dao.ArticleCnt, dao.ArticleAvgWordCnt = oldArticleCnt, oldArticleAvgWordCnt
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return oldArticle, err
	}
	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "ArticleAvgWordCnt").
		Update("value", dao.ArticleAvgWordCnt).Error
	if err != nil {
		tx.Rollback()
		dao.ArticleCnt, dao.ArticleAvgWordCnt = oldArticleCnt, oldArticleAvgWordCnt
		dao.AuthorCnt, dao.AuthorAvgWordCnt = oldAuthorCnt, oldAuthorAvgWordCnt
		return oldArticle, err
	}

	tx.Commit()

	return oldArticle, nil
}

// Target author must be locked before calling this func.
func delAuthor(tx *gorm.DB, id uint64) error {
	var name string
	err := tx.Model(&dao.Author{}).Where("id = ?", id).Select("name").Find(&name).Error
	if err != nil {
		return err
	}
	err = tx.Delete(&dao.Author{ID: id}).Error
	if err != nil {
		return err
	}

	// Del this author id in related inverted indexes.
	nameWords := textToWord(name)

	err = tx.Delete(&dao.AuthorWordCount{ID: id}).Error
	if err != nil {
		return err
	}

	indexes, err := getInvertedIndexesForUpdate(tx, nameWords, word2author)
	if err != nil {
		return err
	}
	for i := range indexes {
		indexes[i].Del(id)
	}
	err = saveInvertedIndexes(tx, nameWords, indexes, word2author)
	if err != nil {
		return err
	}

	// Update global data.
	dao.AuthorCnt--
	dao.AuthorAvgWordCnt = ((float32(dao.AuthorCnt) * dao.AuthorAvgWordCnt) - float32(len(nameWords))) /
		float32(dao.AuthorCnt-1)

	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "AuthorCnt").
		Update("value", dao.AuthorCnt).Error
	if err != nil {
		return err
	}
	err = tx.Model(&dao.Variable{}).Where("`key` = ?", "AuthorAvgWordCnt").
		Update("value", dao.AuthorAvgWordCnt).Error
	if err != nil {
		return err
	}

	return nil
}

func delConnection(tx *gorm.DB, articleID, authorID uint64) error {
	err := tx.Delete(&dao.ArticleToAuthor{ArticleID: articleID, AuthorID: authorID}).Error
	if err != nil {
		return err
	}
	err = tx.Delete(&dao.AuthorToArticle{AuthorID: authorID, ArticleID: articleID}).Error
	if err != nil {
		return err
	}

	return nil
}
