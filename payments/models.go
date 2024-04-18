package main

type CreateAccountReqDto struct {
	AccountId uint32 `json:"account_id"`
	Balance   string `json:"initial_balance"`
}

type GetAccountResDto struct {
	AccountId uint32 `json:"account_id"`
	Balance   string `json:"balance"`
}

type TransactionReqDto struct {
	SourceAccountID      uint32 `json:"source_account_id"`
	DestinationAccountID uint32 `json:"destination_account_id"`
	Amount               string `json:"amount"`
}
