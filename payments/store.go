package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type AccountStore interface {
	GetAccount(ctx context.Context, accountId string) (*Account, error)
	CreateAccount(ctx context.Context, accountId string, balance decimal.Decimal) (*Account, error)
	Transfer(ctx context.Context, fromID, toID string, amount decimal.Decimal) error
}

type SqlStore struct {
	db *pgxpool.Pool
}

var ErrInsufficientFunds = errors.New("insufficient funds in source account")
var ErrAccountNotFound = errors.New("account not found")
var ErrAccountCreation = errors.New("error while account creation")
var ErrAccountRetrieval = errors.New("error while account retrieval")
var ErrAccountNotFoundNoBalance = errors.New("account not found or no balance updated")
var ErrUpdatingAccount = errors.New("error while updating account")
var ErrBeginTx = errors.New("error while begining transaction")
var ErrCommitingTx = errors.New("error while comitting transaction")
var ErrLoggingTx = errors.New("error while logging transaction")

func NewStore(config *Config) *SqlStore {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	log.Println("Database connection successfully established")
	return &SqlStore{db: pool}
}

func (s *SqlStore) GetAccount(ctx context.Context, accountId string) (*Account, error) {
	account := &Account{}
	err := s.db.QueryRow(ctx, "SELECT id, balance FROM accounts WHERE id = $1", accountId).Scan(&account.AccountId, &account.Balance)
	if err != nil {
		log.Printf("Error retrieving account: %v", err)
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (s *SqlStore) CreateAccount(ctx context.Context, accountId string, balance decimal.Decimal) (*Account, error) {
	account := &Account{}
	err := s.db.QueryRow(ctx, "INSERT INTO accounts (id, balance) VALUES ($1, $2) RETURNING id, balance", accountId, balance).Scan(&account.AccountId, &account.Balance)
	if err != nil {
		log.Printf("Error creating account: %v", err)
		return nil, ErrAccountCreation
	}
	return account, nil
}

func (s *SqlStore) Transfer(ctx context.Context, srcAcc, destAcc string, amount decimal.Decimal) error {
	txOptions := pgx.TxOptions{
		IsoLevel: pgx.Serializable, // Specify the isolation level as Serializable
	}
	tx, err := s.db.BeginTx(ctx, txOptions)
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return ErrBeginTx
	}
	defer tx.Rollback(ctx)

	// Determine the order of account IDs to prevent deadlocks
	var firstAcc, secondAcc string
	var firstAmount, secondAmount decimal.Decimal
	if srcAcc < destAcc {
		firstAcc, secondAcc = srcAcc, destAcc
		firstAmount = amount.Neg() // Debit amount (negative)
		secondAmount = amount      // Credit amount (positive)
	} else {
		firstAcc, secondAcc = destAcc, srcAcc
		firstAmount = amount        // Credit amount (positive)
		secondAmount = amount.Neg() // Debit amount (negative)
	}

	log.Printf("Transaction processing from account %s to account %s for amount %s", srcAcc, destAcc, amount)

	// Check balance of the source account to ensure sufficient funds are available
	var balance decimal.Decimal
	err = tx.QueryRow(ctx, "SELECT balance FROM accounts WHERE id = $1 FOR UPDATE", srcAcc).Scan(&balance)
	if err != nil {
		log.Printf("Error retrieving balance for source account: %v", err)
		return ErrAccountRetrieval
	}

	log.Printf("Balance of source account %s is %s", srcAcc, balance)

	if balance.LessThan(amount) {
		log.Printf("Insufficient funds in source account %s of %s", srcAcc, balance)
		return ErrInsufficientFunds
	}

	// Update the balance of the first account
	if cmdTag, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", firstAmount, firstAcc); err != nil {
		log.Printf("Error updating first account: %v", err)
		return ErrUpdatingAccount
	} else if cmdTag.RowsAffected() != 1 {
		log.Printf("First account not found or no balance updated")
		return ErrAccountNotFoundNoBalance
	}

	// Update the balance of the second account
	if cmdTag, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", secondAmount, secondAcc); err != nil {
		log.Printf("Error updating second account: %v", err)
		return ErrUpdatingAccount
	} else if cmdTag.RowsAffected() != 1 {
		log.Printf("Second account not found or no balance updated")
		return ErrAccountNotFoundNoBalance
	}

	// Log the transaction
	if _, err := tx.Exec(ctx, "INSERT INTO transactions (source_account_id, destination_account_id, amount) VALUES ($1, $2, $3)", srcAcc, destAcc, amount); err != nil {
		log.Printf("Error logging transaction: %v", err)
		return ErrLoggingTx
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return ErrCommitingTx
	}

	log.Printf("Transaction successfully processed from account %s to account %s for amount %s", srcAcc, destAcc, amount)
	return nil
}
