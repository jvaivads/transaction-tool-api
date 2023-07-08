package main

import (
	"transaction-tool-api/src/internal/database"
	"transaction-tool-api/src/internal/notifier"
	"transaction-tool-api/src/internal/summarizer"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	controller := summarizer.NewController(
		summarizer.NewService(
			summarizer.NewRepository(
				database.NewSQLClient(
					database.GetLocalMySQLClientConfig(),
				),
			),
			notifier.NewClient(notifier.Options{
				Host:     "smtp.example.com",
				Port:     587,
				Username: "user",
				Password: "password",
			}),
		),
	)

	router.POST("/transaction-tool/resume/:user_id", controller.ResumeTransactions)

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
