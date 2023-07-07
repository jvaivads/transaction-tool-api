package summarizer

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestContext(
	params map[string]string, body [][]string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)

	ctx.Request = &http.Request{}

	for key, value := range params {
		ctx.AddParam(key, value)
	}

	b := &bytes.Buffer{}

	if body != nil {
		w := csv.NewWriter(b)
		if err := w.WriteAll(body); err != nil {
			panic(err)
		}
	}

	ctx.Request.Body = io.NopCloser(b)

	return ctx, r
}

func TestControllerResumeTransactions(t *testing.T) {
	var (
		date = time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)

		bankTransactions = transactions{
			items: []transaction{
				{
					amount: -10,
					date:   date,
				},
				{
					amount: 15,
					date:   date.Add(time.Hour),
				},
			},
			userID: 5,
		}
		customErr = errors.New("custom error")
	)

	tests := []struct {
		name         string
		params       map[string]string
		body         [][]string
		mockApplier  func(m *serviceMock)
		expectedCode int
		expectedBody any
	}{
		{
			name:         "missing user id param",
			expectedCode: http.StatusBadRequest,
			expectedBody: badRequestError("missing user id param"),
		},
		{
			name:         "user id is not an integer",
			params:       map[string]string{"user_id": "bad format"},
			expectedCode: http.StatusBadRequest,
			expectedBody: badRequestError(fmt.Sprintf("user id '%s' is not an integer", "bad format")),
		},
		{
			name:         "csv file is empty",
			params:       map[string]string{"user_id": "5"},
			expectedCode: http.StatusBadRequest,
			expectedBody: badRequestError("csv body is empty"),
		},
		{
			name:   "csv row bad format",
			params: map[string]string{"user_id": "5"},
			body: [][]string{
				{"", "", ""},
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: badRequestError(fmt.Errorf(
				"for row number %d is expected %d elements, however got %d", 1, 2, 3).Error()),
		},
		{
			name:   "service internal error",
			params: map[string]string{"user_id": "5"},
			body: [][]string{
				{"-10", date.Format(time.RFC3339)},
				{"15", date.Add(time.Hour).Format(time.RFC3339)},
			},
			mockApplier: func(m *serviceMock) {
				m.On("notifyResume", bankTransactions).Return(customErr).Once()
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: Error{
				Status:  http.StatusInternalServerError,
				Code:    "internal error",
				Message: customErr.Error(),
			},
		},
		{
			name:   "service internal error",
			params: map[string]string{"user_id": "5"},
			body: [][]string{
				{"-10", date.Format(time.RFC3339)},
				{"15", date.Add(time.Hour).Format(time.RFC3339)},
			},
			mockApplier: func(m *serviceMock) {
				m.On("notifyResume", bankTransactions).Return(nil).Once()
			},
			expectedCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				ctx, r   = getTestContext(test.params, test.body)
				servMock = &serviceMock{}
				ctl      = controller{service: servMock}
			)

			if test.mockApplier != nil {
				test.mockApplier(servMock)
				defer servMock.AssertExpectations(t)
			}

			expectedBody := ""
			if test.expectedBody != nil {
				b, err := json.Marshal(test.expectedBody)
				require.Nil(t, err)
				expectedBody = string(b)
			}

			ctl.ResumeTransactions(ctx)

			assert.Equal(t, test.expectedCode, r.Code)
			assert.Equal(t, expectedBody, r.Body.String())
		})
	}
}

func TestControllerParseFileToTransactions(t *testing.T) {
	var (
		date = time.Date(2021, 10, 1, 0, 0, 0, 0, time.UTC)

		bankTransactions = transactions{
			items: []transaction{
				{
					amount: -10.5,
					date:   date,
				},
				{
					amount: 15,
					date:   date.Add(time.Hour),
				},
			},
			userID: 5,
		}
	)

	tests := []struct {
		name           string
		transactions   [][]string
		expectedResult transactions
		expectedErr    error
	}{
		{
			name:           "no transactions",
			transactions:   nil,
			expectedResult: transactions{userID: 5, items: make([]transaction, 0)},
			expectedErr:    nil,
		},
		{
			name: "elements by row unexpected",
			transactions: [][]string{
				{""},
			},
			expectedErr: fmt.Errorf(
				"for row number %d is expected %d elements, however got %d", 1, 2, 1),
		},
		{
			name: "amount bad format",
			transactions: [][]string{
				{"bad format", ""},
			},
			expectedErr: fmt.Errorf(
				"error parsing float amount (%s) from row number %d", "bad format", 1),
		},
		{
			name: "amount zero no valid",
			transactions: [][]string{
				{"0", ""},
			},
			expectedErr: fmt.Errorf(
				"for row number %d, transaction amount is zero", 1),
		},
		{
			name: "date bad format",
			transactions: [][]string{
				{"-10.5", "bad format"},
			},
			expectedErr: fmt.Errorf(
				"error parsing date '%s' because of no compliance with RFC3339 layout for row number %d",
				"bad format", 1),
		},
		{
			name: "successfully parsing",
			transactions: [][]string{
				{"-10.5", date.Format(time.RFC3339)},
				{"15", date.Add(time.Hour).Format(time.RFC3339)},
			},
			expectedResult: bankTransactions,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctl := controller{}

			result, err := ctl.parseFileToTransactions(test.transactions, 5)

			assert.Equal(t, test.expectedResult, result)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}
