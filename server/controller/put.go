package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/service"
)

func Update(c *gin.Context) {
	tmpArticle := service.Article{}
	err := c.Bind(&tmpArticle)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newArticle, err := service.Update(tmpArticle)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, newArticle)
	}
}
