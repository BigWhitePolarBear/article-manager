package main

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	reAuthor        = regexp.MustCompile(".*<author.*>(.*?)</author>.*")
	reTitle         = regexp.MustCompile(".*<title>(.*?)</title>.*")
	reBookTitle     = regexp.MustCompile(".*<booktitle>(.*?)</booktitle>.*")
	reJournal       = regexp.MustCompile(".*<journal>(.*?)</journal>.*")
	reVolume        = regexp.MustCompile(".*<volume>(.*?)</volume>.*")
	rePages         = regexp.MustCompile(".*<pages>(.*?)</pages>.*")
	reYear          = regexp.MustCompile(".*<year>(.*?)</year>.*")
	reEE            = regexp.MustCompile(".*<ee.*>(.*?)</ee>.*")
	reArticle       = regexp.MustCompile(".*<article.*")
	reBook          = regexp.MustCompile(".*<book.*")
	reCollection    = regexp.MustCompile(".*<incollection.*")
	reProceeding    = regexp.MustCompile(".*<inproceedings.*")
	reArticleEnd    = regexp.MustCompile(".*</article>.*")
	reBookEnd       = regexp.MustCompile(".*</book>.*")
	reCollectionEnd = regexp.MustCompile(".*</incollection>.*")
	reProceedingEnd = regexp.MustCompile(".*</inproceedings>.*")
)

// Delete the number at the end of names.
func repairName(name string) string {
	words := strings.Fields(name)
	nameBuilder := strings.Builder{}
	for _, word := range words {
		if strings.ContainsAny(word, "0123456789") {
			continue
		}
		if nameBuilder.Len() == 0 {
			nameBuilder.WriteString(word)
		} else {
			nameBuilder.WriteByte(' ')
			nameBuilder.WriteString(word)
		}
	}
	return nameBuilder.String()
}

func match() {
	defer wg.Done()
	// Tags for speeding up regexp matches.
	inArticle, inAuthor, inEE, authorDone, titleDone, bookDone, journalDone, volumeDone, pagesDone, yearDone, eeDone :=
		false, false, false, false, false, false, false, false, false, false, false

	article := Article{}
	authors := ""

Loop:
	for {
		select {
		case line := <-read2match:
			if inArticle && (reBookEnd.Match(line) || reArticleEnd.Match(line) ||
				reProceedingEnd.Match(line) || reCollectionEnd.Match(line)) {
				inArticle = false
				if article.Title == "" {
					continue Loop
				}
				articleMatch2Record <- article
				authorsMatch2Record <- authors
				continue Loop
			}

			if !inArticle && (reBook.Match(line) || reArticle.Match(line) ||
				reProceeding.Match(line) || reCollection.Match(line)) {
				inArticle = true
				article = Article{}
				authors = ""
				authorDone, titleDone, bookDone, journalDone, volumeDone, pagesDone, yearDone, eeDone =
					false, false, false, false, false, false, false, false
				continue Loop
			}

			if inArticle {
				if !authorDone {
					temp := reAuthor.FindSubmatch(line)
					if temp != nil {
						inAuthor = true
						name := temp[1]
						t := repairName(string(name))
						if authors == "" {
							authors = t
						} else {
							authors += ", " + t
						}
						continue Loop
					} else if inAuthor {
						inAuthor = false
						authorDone = true
					}
				}
				if !titleDone {
					temp := reTitle.FindSubmatch(line)
					if temp != nil {
						title := temp[1]
						article.Title = string(title)
						titleDone = true
						continue Loop
					}
				}
				if !bookDone {
					temp := reBookTitle.FindSubmatch(line)
					if temp != nil {
						bookTitle := string(temp[1])
						if len(bookTitle) == 0 {
							continue Loop
						}
						var id uint64
						DB.Model(&Book{}).Select("id").Where("name = ?", bookTitle).Find(&id)
						if id == 0 {
							newBook := Book{Name: bookTitle}
							DB.Model(&Book{}).Create(&newBook)
							article.BookID = &newBook.ID
						} else {
							article.BookID = &id
						}
						bookDone = true
						continue Loop
					}
				}
				if !journalDone {
					temp := reJournal.FindSubmatch(line)
					if temp != nil {
						journal := string(temp[1])
						if len(journal) == 0 {
							continue Loop
						}
						var id uint64
						DB.Model(&Journal{}).Select("id").Where("name = ?", journal).Find(&id)
						if id == 0 {
							newJournal := Journal{Name: journal}
							DB.Model(&Journal{}).Create(&newJournal)
							article.JournalID = &newJournal.ID
						} else {
							article.JournalID = &id
						}
						journalDone = true
						continue Loop
					}
				}
				if !volumeDone {
					temp := reVolume.FindSubmatch(line)
					if temp != nil {
						volume := temp[1]
						article.Volume = string(volume)
						volumeDone = true
						continue Loop
					}
				}
				if !pagesDone {
					temp := rePages.FindSubmatch(line)
					if temp != nil {
						pages := temp[1]
						article.Pages = string(pages)
						pagesDone = true
						continue Loop
					}
				}
				if !yearDone {
					temp := reYear.FindSubmatch(line)
					if temp != nil {
						year := temp[1]
						t, _ := strconv.ParseUint(string(year), 10, 16)
						if t == 0 {
							continue
						}
						t_ := uint16(t)
						article.Year = &t_
						yearDone = true
						continue Loop
					}
				}
				if !eeDone {
					temp := reEE.FindSubmatch(line)
					if temp != nil {
						inEE = true
						ee := temp[1]
						if article.EE == "" {
							article.EE = string(ee)
						} else {
							article.EE += ", " + string(ee)
						}
					} else if inEE {
						inEE = false
						eeDone = true
					}
				}
			}
		case <-readOK:
			matchOK <- struct{}{}
		}
	}
}
