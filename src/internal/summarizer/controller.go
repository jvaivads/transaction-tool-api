package summarizer

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type controller struct {
	service Service
}

type Error struct {
	Status  int
	Code    string
	Message string
}

func badRequestError(message string) Error {
	return Error{
		Status:  http.StatusBadRequest,
		Code:    "bad request",
		Message: message,
	}
}

func (ctl controller) ResumeTransactions(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, badRequestError("missing user id param"))
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, badRequestError(fmt.Sprintf("user id '%s' is not an integer", userIDStr)))
		return
	}

	reader := csv.NewReader(c.Request.Body)
	file, err := reader.ReadAll()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			badRequestError(fmt.Sprintf("error reading csv body due to: %s", err.Error())))
		return
	}

	bankTransactions, err := ctl.parseFileToTransactions(file, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, badRequestError(err.Error()))
		return
	}

	if len(bankTransactions.items) == 0 {
		c.JSON(http.StatusBadRequest, badRequestError("csv body is empty"))
		return
	}

	if err = ctl.service.notifyResume(c.Request.Context(), bankTransactions); err != nil {
		c.JSON(http.StatusInternalServerError, Error{
			Status:  http.StatusInternalServerError,
			Code:    "internal error",
			Message: err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}

func (ctl controller) parseFileToTransactions(file [][]string, userID int64) (transactions, error) {
	var (
		totalElementsByRow = 2
		bankTransactions   = transactions{userID: userID, items: make([]transaction, 0, len(file))}
	)

	for i, row := range file {
		if len(row) != totalElementsByRow {
			return transactions{}, fmt.Errorf(
				"for row number %d is expected %d elements, however got %d", i+1, totalElementsByRow, len(row))
		}

		amount, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			return transactions{}, fmt.Errorf(
				"error parsing float amount (%s) from row number %d", row[0], i+1)
		}

		if amount == 0 {
			return transactions{}, fmt.Errorf(
				"for row number %d, transaction amount is zero", i+1)
		}

		date, err := time.Parse(time.RFC3339, row[1])
		if err != nil {
			return transactions{}, fmt.Errorf(
				"error parsing date '%s' because of no compliance with RFC3339 layout for row number %d",
				row[1], i+1)
		}

		bankTransactions.items = append(bankTransactions.items, transaction{
			amount: amount,
			date:   date,
		})
	}
	return bankTransactions, nil
}
