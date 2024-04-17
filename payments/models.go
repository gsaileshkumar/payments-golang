package main

type CreateAccountReqDto struct {
	AccountId string `json:"account_id"`
	Balance   string `json:"initial_balance"`
}

type GetAccountResDto struct {
	AccountId string `json:"account_id"`
	Balance   string `json:"balance"`
}

type TransactionReqDto struct {
	SourceAccountID      string `json:"source_account_id"`
	DestinationAccountID string `json:"destination_account_id"`
	Amount               string `json:"amount"`
}
