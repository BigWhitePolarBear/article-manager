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
	//c.GET("/authors/search", controller.GetAuthor)
	//c.GET("/authors/top/:page", controller.GetTopAuthor)

	// Use middleware to check if client has root privilege.
	adminRoutine := c.Group("/admin", controller.LoginMidWare())
	{
		// Root retrieve to get work id.
		adminRoutine.GET("/papers/search", controller.AdminMark, controller.SearchArticle)

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
