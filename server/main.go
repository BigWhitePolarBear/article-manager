package main

import (
	"github.com/gin-gonic/gin"
	"server/controller"
	"server/dao"
	"server/service"
)

func main() {
	dao.Init()
	service.Init()

	c := gin.Default()

	// Retrieve
	c.GET("/article/retrieve", controller.SearchArticle)
	c.GET("/author/retrieve", controller.SearchAuthor)
	c.GET("/authors/top/:page", controller.GetTopAuthor)

	// Use middleware to check if client has root privilege.
	adminRoutine := c.Group("/admin", controller.AdminMark) //, controller.LoginMidWare())
	{
		// Root retrieve to get article id.
		adminRoutine.GET("/article/retrieve", controller.SearchArticle)
		adminRoutine.GET("/author/retrieve", controller.SearchAuthor)
		adminRoutine.GET("/authors/top/:page", controller.GetTopAuthor)

		// Only support moderation on article data.

		// Create
		adminRoutine.POST("/create", controller.Post)
		//
		//// Update
		//adminRoutine.PUT("/update", controller.UpdateArticle)
		//
		// Delete
		adminRoutine.DELETE("/delete", controller.Del)

	}

	err := c.Run(":80")
	if err != nil {
		panic(err)
	}
}
