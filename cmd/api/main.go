module cc.io/arena

// main.go
package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	// Setup routes, middleware, etc.
	app.Run() // Run the server
}