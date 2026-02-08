// reponsible forrunning database
package main

import "github.com/gin-gonic/gin"

func main() {
	// create server, take a look at routes, want api fast, use instance from the memory, pointer variable
	// * is a pointer, reference something in the memory
	// pointer refers to the address or instance in memory, and not copy entire thing
	var router *gin.Engine = gin.Default() // gin => do client request and response
	router.GET("/", func(c *gin.Context) {
		// gin.H is a shortcut for map[string]interface{} or map[string]any
		c.JSON(200, gin.H{
			"message": "todo api running successfully!!!",
			"status":  "success",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
