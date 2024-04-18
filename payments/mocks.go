package main

import (
	"context"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

type MockAccountStore struct {
	mock.Mock
}

func (m *MockAccountStore) GetAccount(ctx context.Context, accountId uint32) (*Account, error) {
	args := m.Called(ctx, accountId)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockAccountStore) CreateAccount(ctx context.Context, accountId uint32, balance decimal.Decimal) (*Account, error) {
	args := m.Called(ctx, accountId, balance)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockAccountStore) Transfer(ctx context.Context, srcAccId, destAccId uint32, amount decimal.Decimal) error {
	args := m.Called(ctx, srcAccId, destAccId, amount)
	return args.Error(0)
}
