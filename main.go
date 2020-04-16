package main

import (
	"github.com/gin-gonic/gin"
	"todo.com/controllers"
)

func main() {

	router := gin.Default()

	c := new(controllers.TodoController)

	todo := router.Group("/todo")
	{
		todo.GET("/readAll", c.Read)
		todo.GET("/read/:id", c.ReadById)
		todo.POST("/create", c.Create)
		todo.GET("/remove/:id", c.Remove)
		todo.POST("/update", c.Update)
	}

	router.Run()

}
