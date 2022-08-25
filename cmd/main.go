package main

import (
	"log"
	"os"
	"wxrobot/worker"

	"github.com/gin-gonic/gin"
)

func init() {
	os.Setenv("WECHATY_PUPPET_SERVICE_TOKEN", "token")
}

func main() {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		log.Printf("uri: %v ,body: %v\n", c.Request.URL.String(), c.Request.Body)
	})
	{

	}
	go worker.Worker()
	router.Run(":8080")
}
