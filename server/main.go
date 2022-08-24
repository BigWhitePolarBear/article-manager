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
	c.GET("/article/search", controller.SearchArticle)
	c.GET("/author/search", controller.SearchAuthor)
	c.GET("/authors/top/:page", controller.GetTopAuthor)

	// Use middleware to check if client has root privilege.
	adminRoutine := c.Group("/admin", controller.AdminMark) //, controller.LoginMidWare())
	{
		// Root retrieve to get work id.
		adminRoutine.GET("/article/search", controller.SearchArticle)
		adminRoutine.GET("/author/search", controller.SearchAuthor)
		adminRoutine.GET("/authors/top/:page", controller.GetTopAuthor)

		//// Create
		//adminRoutine.POST("/add", post.Work)
		//
		//// Update
		//adminRoutine.PUT("/update", controller.UpdateArticle)
		//
		//// Delete
		//adminRoutine.DELETE("/delete", controller.DelWork)
	}

	err := c.Run(":80")
	if err != nil {
		panic(err)
	}
}
