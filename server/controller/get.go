package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/service"
	"strconv"
)

func SearchArticle(c *gin.Context) {
	queries := make(map[service.QueryType]string)
	if title, ok := c.GetQuery("title"); ok {
		queries[service.TitleQuery] = title
	} else {
		c.JSON(http.StatusBadRequest, "Must input title queries!")
	}

	if author, ok := c.GetQuery("author"); ok {
		queries[service.AuthorQuery] = author
	}

	if not, ok := c.GetQuery("not"); ok {
		queries[service.NotQuery] = not
	}

	page := c.DefaultQuery("page", "1")
	if page == "0" {
		c.JSON(http.StatusBadRequest, "There is no 0th page!")
	}
	queries[service.PageQuery] = page

	_, admin := c.Get("admin")
	results, err := service.SearchArticle(queries, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, results)
	}
}

func SearchAuthor(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	if page == "0" {
		c.JSON(http.StatusBadRequest, "There is no 0th page!")
	}

	name, ok := c.GetQuery("name")
	if !ok {
		c.JSON(http.StatusBadRequest, "Please enter name!")
		return
	}

	_, admin := c.Get("admin")
	results, err := service.SearchAuthor(name, page, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, results)
	}
}

func GetTopAuthor(c *gin.Context) {
	page, err := strconv.ParseUint(c.Param("page"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	_, admin := c.Get("admin")
	results, err := service.GetTopAuthor(page, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, results)
	}
}
