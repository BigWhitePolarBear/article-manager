package service

import (
	"math"
	"paperSearchServer/dao"
	"strings"
)

type Article struct {
	ID        *uint `json:",omitempty"`
	Title     string
	Authors   string  `json:",omitempty"`
	Journal   string  `json:",omitempty"`
	Volume    string  `json:",omitempty"`
	Month     string  `json:",omitempty"`
	Year      *uint16 `json:",omitempty"`
	CdRom     string  `json:",omitempty"`
	EE        string  `json:",omitempty"`
	Publisher string  `json:",omitempty"`
	ISBN      string  `json:",omitempty"`
}

type Author struct {
	Name         string
	Works        string `json:",omitempty"`
	ArticleCount uint16
}

type QueryType uint

const (
	TitleQuery = iota
	YearQuery
	AuthorQuery
)

func GetWork(query map[QueryType]string, limit, offset int, admin bool) (*[]Article, error) {
	searchSQL := worksTable.Limit(limit).Offset(offset)
	if query[TitleQuery] != "" {
		searchSQL = searchSQL.Where("title = ?", query[TitleQuery])
	}
	if query[YearQuery] != "" {
		searchSQL = searchSQL.Where("year = ?", query[YearQuery])
	}
	if query[AuthorQuery] != "" {
		names := strings.Split(query[AuthorQuery], ",")
		// Temporarily stored the ids.
		var idLists []string
		authorsTable.Where("name in ?", names).Select("works").Find(&idLists)

		// Record each set of work id, for finding out their intersection.
		// 'l' means length.
		resL := math.MaxInt
		tempMaps := make([]map[string]struct{}, 0, len(idLists))
		for _, _idList := range idLists {
			l := len(_idList)
			if l < resL {
				resL = l
			}
			tempMap := make(map[string]struct{}, l)
			idList := strings.Fields(_idList)
			for _, id := range idList {
				tempMap[id] = struct{}{}
			}
			tempMaps = append(tempMaps, tempMap)
		}
		workIDs := make([]string, 0, resL)
		// Drop the ids those aren't exist in the work list of every author.
	Loop:
		for id := range tempMaps[0] {
			for _, M := range tempMaps {
				if _, ok := M[id]; !ok {
					continue Loop
				}
			}
			workIDs = append(workIDs, id)
		}
		searchSQL = searchSQL.Where("id in ?", workIDs)
	}

	_results := []dao.Article{}
	err := searchSQL.Find(&_results).Error
	if err != nil {
		return nil, err
	} else {
		results := make([]Article, 0, len(_results))
		for _, _result := range _results {
			temp := Article{}
			if admin {
				temp.ID = &_result.ID
			}
			temp.Title, temp.ISBN, temp.CdRom, temp.Month, temp.Year, temp.Volume, temp.EE =
				_result.Title, _result.ISBN, _result.CdRom, _result.Month, _result.Year, _result.Volume, _result.EE

			if _result.JournalID != nil {
				journalsTable.Where("id = ?", *_result.JournalID).Select("name").Find(&temp.Journal)
			}
			if _result.PublisherID != nil {
				publishersTable.Where("id = ?", *_result.PublisherID).Select("name").Find(&temp.Publisher)
			}

			authorID := strings.Fields(_result.Authors)
			authorBuilder := strings.Builder{}
			for i, id := range authorID {
				if i > 0 {
					authorBuilder.WriteString(", ")
				}
				name := ""
				authorsTable.Where("id = ?", id).Select("name").Find(&name)
				authorBuilder.WriteString(name)
			}
			temp.Authors = authorBuilder.String()
			results = append(results, temp)
		}
		return &results, nil
	}
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
