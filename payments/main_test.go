package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing" // Assuming this is the decimal package used

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAccountHandler(t *testing.T) {
	mockStore := new(MockAccountStore)
	app := &Payments{store: mockStore}

	mockStore.On("GetAccount", mock.Anything, mock.Anything).Return(&Account{}, ErrAccountNotFound)
	mockStore.On("CreateAccount", mock.Anything, mock.Anything, mock.Anything).Return(&Account{AccountId: 12345, Balance: "100.50"}, nil)

	testServer := httptest.NewServer(http.HandlerFunc(app.CreateAccountHandler))
	defer testServer.Close()

	account := CreateAccountReqDto{
		AccountId: 12345,
		Balance:   "100.50",
	}

	body, _ := json.Marshal(account)
	request, _ := http.NewRequest("POST", testServer.URL+"/accounts", bytes.NewBuffer(body))
	response, err := http.DefaultClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.StatusCode)

	mockStore.AssertExpectations(t)

}

func TestCreateAccountHandler_AccountExists(t *testing.T) {
	mockStore := new(MockAccountStore)
	app := &Payments{store: mockStore}

	// Mock GetAccount to simulate an existing account
	mockStore.On("GetAccount", mock.Anything, uint32(12345)).Return(&Account{AccountId: 12345, Balance: "100.50"}, nil)

	testServer := httptest.NewServer(http.HandlerFunc(app.CreateAccountHandler))
	defer testServer.Close()

	account := CreateAccountReqDto{
		AccountId: 12345,
		Balance:   "100.50",
	}

	body, _ := json.Marshal(account)
	request, _ := http.NewRequest("POST", testServer.URL+"/accounts", bytes.NewBuffer(body))
	response, err := http.DefaultClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	mockStore.AssertExpectations(t)
}

func TestCreateAccountHandler_FailureCreatingAccount(t *testing.T) {
	mockStore := new(MockAccountStore)
	app := &Payments{store: mockStore}

	// Assume account does not exist
	mockStore.On("GetAccount", mock.Anything, uint32(12345)).Return(&Account{}, ErrAccountNotFound)
	// Simulate failure during account creation
	mockStore.On("CreateAccount", mock.Anything, uint32(12345), mock.Anything).Return(&Account{}, errors.New("failed to create account"))

	testServer := httptest.NewServer(http.HandlerFunc(app.CreateAccountHandler))
	defer testServer.Close()

	account := CreateAccountReqDto{
		AccountId: 12345,
		Balance:   "100.50",
	}

	body, _ := json.Marshal(account)
	request, _ := http.NewRequest("POST", testServer.URL+"/accounts", bytes.NewBuffer(body))
	response, err := http.DefaultClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode)

	mockStore.AssertExpectations(t)
}

func TestCreateAccountHandler_BadRequest(t *testing.T) {
	app := &Payments{store: nil} // Store not needed for this test

	testServer := httptest.NewServer(http.HandlerFunc(app.CreateAccountHandler))
	defer testServer.Close()

	// Send invalid JSON
	body := bytes.NewBufferString("{invalidJson:}")
	request, _ := http.NewRequest("POST", testServer.URL+"/accounts", body)
	response, err := http.DefaultClient.Do(request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestTransferAmountHandler_Success(t *testing.T) {
	mockStore := new(MockAccountStore)
	app := &Payments{store: mockStore}

	mockStore.On("Transfer", mock.Anything, uint32(1), uint32(2), mock.Anything).Return(nil)

	transaction := TransactionReqDto{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "100.00",
	}
	body, _ := json.Marshal(transaction)
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.TransferAmountHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockStore.AssertExpectations(t)
}

func TestTransferAmountHandler_InvalidInput(t *testing.T) {
	app := &Payments{store: nil} // Store not involved in this test

	// Testing various invalid scenarios
	tests := []struct {
		name     string
		payload  TransactionReqDto
		expected int
	}{
		{"Same Account IDs", TransactionReqDto{1, 1, "50.00"}, http.StatusBadRequest},
		{"Zero Source Account ID", TransactionReqDto{0, 2, "20.00"}, http.StatusBadRequest},
		{"Zero Destination Account ID", TransactionReqDto{1, 0, "15.00"}, http.StatusBadRequest},
		{"Empty Amount", TransactionReqDto{1, 2, ""}, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(app.TransferAmountHandler)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tc.expected, rr.Code)
		})
	}
}

func TestTransferAmountHandler_DbError(t *testing.T) {
	mockStore := new(MockAccountStore)
	app := &Payments{store: mockStore}

	mockStore.On("Transfer", mock.Anything, uint32(1), uint32(2), mock.Anything).Return(errors.New("db error"))

	transaction := TransactionReqDto{
		SourceAccountID:      1,
		DestinationAccountID: 2,
		Amount:               "100.00",
	}
	body, _ := json.Marshal(transaction)
	req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.TransferAmountHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockStore.AssertExpectations(t)
}

func TestCatchAllHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/nonexistentpath", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(catchAllHandler)
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Check the response body is what we expect.
	expected := map[string]string{"error": "Route not found"}
	var actual map[string]string
	err := json.NewDecoder(rr.Body).Decode(&actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}
