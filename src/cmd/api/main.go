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
			notifier.NewClient(notifier.GetOptions()),
		),
	)

	router.POST("/transaction-tool/resume/:user_id", controller.ResumeTransactions)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
