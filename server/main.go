package main

import (
	"github.com/gin-gonic/gin"
	"server/controller"
)

func main() {
	c := gin.Default()

	// Retrieve
	c.GET("/article/search", controller.SearchArticle)
	c.GET("/authors/search", controller.GetAuthor)
	c.GET("/authors/top/:page", controller.GetTopAuthor)

	// Use middleware to check if client has root privilege.
	adminRoutine := c.Group("/admin", controller.LoginMidWare())
	{
		// Root retrieve to get work id.
		adminRoutine.GET("/papers/search", controller.AdminMark, controller.GetWork)

		// Create
		adminRoutine.POST("/add", post.Work)

		// Update
		adminRoutine.PUT("/update", controller.UpdateArticle)

		// Delete
		adminRoutine.DELETE("/delete", controller.DelWork)
	}

	c.Run(":8000")
}
