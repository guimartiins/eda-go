package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidAmount = errors.New("invalid amount")
var ErrInsufficientFunds = errors.New("insufficient funds")

type Transaction struct {
	ID          string
	AccountFrom *Account
	AccountTo   *Account
	Amount      float64
	CreatedAt   time.Time
}

func NewTransaction(accountFrom *Account, accountTo *Account, amount float64) (*Transaction, error) {
	transaction := &Transaction{
		ID:          uuid.New().String(),
		AccountFrom: accountFrom,
		AccountTo:   accountTo,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}

	if err := transaction.Validate(); err != nil {
		return nil, err
	}

	transaction.Commit()
	return transaction, nil
}

func (t *Transaction) Validate() error {
	if t.Amount <= 0 {
		return ErrInvalidAmount
	}

	if t.AccountFrom.Balance < t.Amount {
		return ErrInsufficientFunds
	}
	return nil
}

func (t *Transaction) Commit() {
	t.AccountFrom.Debit(t.Amount)
	t.AccountTo.Credit(t.Amount)
}
