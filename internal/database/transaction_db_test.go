package database

import (
	"database/sql"
	"testing"

	"github.com/guimartiins/eda-go/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type TransactionDBTestSuite struct {
	suite.Suite
	db            *sql.DB
	transactionDB *TransactionDB
	client1       *entity.Client
	client2       *entity.Client
	account1      *entity.Account
	account2      *entity.Account
}

func (s *TransactionDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	s.db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at, date updated_at date)")
	s.db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	s.db.Exec("CREATE TABLE transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount float, created_at date)")
	s.transactionDB = NewTransactionDB(db)
	s.client1, _ = entity.NewClient("John", "j@j.com")
	s.client2, _ = entity.NewClient("Jane", "j2@j.com")
	s.account1 = entity.NewAccount(s.client1)
	s.account1.Credit(1000)
	s.account2 = entity.NewAccount(s.client2)
	s.account2.Credit(1000)

	// Insert test clients and accounts
	s.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		s.client1.ID, s.client1.Name, s.client1.Email, s.client1.CreatedAt)
	s.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		s.client2.ID, s.client2.Name, s.client2.Email, s.client2.CreatedAt)
	s.db.Exec("INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)",
		s.account1.ID, s.account1.Client.ID, s.account1.Balance, s.account1.CreatedAt)
	s.db.Exec("INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)",
		s.account2.ID, s.account2.Client.ID, s.account2.Balance, s.account2.CreatedAt)
}

func (s *TransactionDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE transactions")
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func (s *TransactionDBTestSuite) TestCreate() {
	transaction, err := entity.NewTransaction(s.account1, s.account2, 100)
	s.Nil(err)
	err = s.transactionDB.Create(transaction)
	s.Nil(err)

	var savedTransaction entity.Transaction
	savedTransaction.AccountFrom = &entity.Account{}
	savedTransaction.AccountTo = &entity.Account{}

	row := s.db.QueryRow("SELECT id, account_id_from, account_id_to, amount FROM transactions WHERE id = ?", transaction.ID)
	err = row.Scan(
		&savedTransaction.ID,
		&savedTransaction.AccountFrom.ID,
		&savedTransaction.AccountTo.ID,
		&savedTransaction.Amount,
	)

	s.Nil(err)
	s.Equal(transaction.ID, savedTransaction.ID)
	s.Equal(transaction.AccountFrom.ID, savedTransaction.AccountFrom.ID)
	s.Equal(transaction.AccountTo.ID, savedTransaction.AccountTo.ID)
	s.Equal(transaction.Amount, savedTransaction.Amount)
}

func TestTransactionDBTestSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTestSuite))
}
