package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
