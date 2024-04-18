package main

import (
	"context"

	"github.com/shopspring/decimal"
)

type Account struct {
	AccountId string
	Balance   string
}

type AccountStore interface {
	GetAccount(ctx context.Context, accountId string) (*Account, error)
	CreateAccount(ctx context.Context, accountId string, balance decimal.Decimal) (*Account, error)
	Transfer(ctx context.Context, fromID, toID string, amount decimal.Decimal) error
}
