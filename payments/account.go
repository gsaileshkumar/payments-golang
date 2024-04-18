package main

import (
	"context"

	"github.com/shopspring/decimal"
)

type Account struct {
	AccountId uint32
	Balance   string
}

type AccountStore interface {
	GetAccount(ctx context.Context, accountId uint32) (*Account, error)
	CreateAccount(ctx context.Context, accountId uint32, balance decimal.Decimal) (*Account, error)
	Transfer(ctx context.Context, srcAccId, destAccId uint32, amount decimal.Decimal) error
}
