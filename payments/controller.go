package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

var ErrGeneric = errors.New("server error")
var ErrBadRequest = errors.New("bad request")
var ErrAccExists = errors.New("account already exists")

func sendJSONError(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonResp := map[string]string{"error": err.Error()}
	json.NewEncoder(w).Encode(jsonResp)
}

func (p *Payments) CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateAccountHandler")
	var account CreateAccountReqDto
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		sendJSONError(w, ErrBadRequest, http.StatusBadRequest)
		return
	}

	if account.AccountId == 0 || account.Balance == "" {
		sendJSONError(w, ErrBadRequest, http.StatusBadRequest)
		return
	}
	_, err = p.store.GetAccount(r.Context(), account.AccountId)
	if err == nil {
		sendJSONError(w, ErrAccExists, http.StatusBadRequest)
		return
	}
	balance, err := decimal.NewFromString(account.Balance)
	if err != nil {
		sendJSONError(w, ErrGeneric, http.StatusInternalServerError)
		return
	}

	acc, err := p.store.CreateAccount(r.Context(), account.AccountId, balance)
	if err != nil {
		sendJSONError(w, ErrGeneric, http.StatusInternalServerError)
		return
	}

	respBytes, err := json.Marshal(&GetAccountResDto{
		AccountId: acc.AccountId,
		Balance:   acc.Balance,
	})
	if err != nil {
		sendJSONError(w, ErrGeneric, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(respBytes)
}

func (p *Payments) GetAccounDetailsHandler(w http.ResponseWriter, r *http.Request) {
	accountId := r.PathValue("id")

	log.Println("GetAccounDetailsHandler for account: ", accountId)

	if accountId == "" {
		sendJSONError(w, ErrBadRequest, http.StatusBadRequest)
		return
	}

	accountIdInt, err := strconv.ParseUint(accountId, 10, 32)
	if err != nil {
		sendJSONError(w, ErrBadRequest, http.StatusBadRequest)
		return
	}

	account, err := p.store.GetAccount(r.Context(), uint32(accountIdInt))
	if err != nil {
		sendJSONError(w, err, http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(&GetAccountResDto{
		AccountId: account.AccountId,
		Balance:   account.Balance,
	})
}

func (p *Payments) TransferAmountHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("TransferAmountHandler")
	var transaction TransactionReqDto
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		sendJSONError(w, ErrBadRequest, http.StatusBadRequest)
		return
	}

	if transaction.SourceAccountID == transaction.DestinationAccountID ||
		transaction.SourceAccountID == 0 || transaction.DestinationAccountID == 0 ||
		transaction.Amount == "" {
		sendJSONError(w, ErrBadRequest, http.StatusBadRequest)
		return
	}

	amount, err := decimal.NewFromString(transaction.Amount)
	if err != nil {
		sendJSONError(w, ErrGeneric, http.StatusInternalServerError)
		return
	}

	// Log the transaction in the database for audit
	err = p.store.Transfer(r.Context(), transaction.SourceAccountID, transaction.DestinationAccountID, amount)
	if err != nil {
		sendJSONError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
