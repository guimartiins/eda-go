package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {
	client1, _ := NewClient("John", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane", "j@j2.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 100)
	assert.Nil(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 900.0, account1.Balance)
	assert.Equal(t, 1100.0, account2.Balance)
}

func TestCreateTransactionWithInsufficientFunds(t *testing.T) {
	client1, _ := NewClient("John", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane", "j@j2.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 2000)
	assert.NotNil(t, err)
	assert.Error(t, err, ErrInsufficientFunds)
	assert.Nil(t, transaction)
	assert.Equal(t, 1000.0, account1.Balance)
	assert.Equal(t, 1000.0, account2.Balance)
}

func TestCreateTransactionWithInvalidAmount(t *testing.T) {
	client1, _ := NewClient("John", "j@j.com")
	account1 := NewAccount(client1)
	client2, _ := NewClient("Jane", "j@j2.com")
	account2 := NewAccount(client2)

	account1.Credit(1000)
	account2.Credit(1000)

	transaction, err := NewTransaction(account1, account2, 0)
	assert.NotNil(t, err)
	assert.Error(t, err, ErrInvalidAmount)
	assert.Nil(t, transaction)
	assert.Equal(t, 1000.0, account1.Balance)
	assert.Equal(t, 1000.0, account2.Balance)
}
