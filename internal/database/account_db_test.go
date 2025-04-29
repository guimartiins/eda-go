package database

import (
	"database/sql"
	"testing"

	"github.com/guimartiins/fcutils/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type AccountDBTestSuite struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
	client    *entity.Client
}

func (s *AccountDBTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	s.db.Exec("CREATE TABLE clients (id varchar(255), name varchar(255), email varchar(255), created_at date)")
	s.db.Exec("CREATE TABLE accounts (id varchar(255), client_id varchar(255), balance float, created_at date)")
	s.accountDB = NewAccountDB(db)
	s.client, _ = entity.NewClient("John", "j@j.com")
	s.db.Exec("INSERT INTO clients (id, name, email, created_at) VALUES (?, ?, ?, ?)",
		s.client.ID, s.client.Name, s.client.Email, s.client.CreatedAt)
}

func (s *AccountDBTestSuite) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("DROP TABLE accounts")
	s.db.Exec("DROP TABLE clients")
}

func (s *AccountDBTestSuite) TestSave() {
	account := entity.NewAccount(s.client)
	err := s.accountDB.Save(account)
	s.Nil(err)
}

func (s *AccountDBTestSuite) TestGet() {
	account := entity.NewAccount(s.client)
	account.Credit(100)
	s.accountDB.Save(account)

	accountDB, err := s.accountDB.Get(account.ID)
	s.Nil(err)
	s.Equal(account.ID, accountDB.ID)
	s.Equal(account.Client.ID, accountDB.Client.ID)
	s.Equal(account.Client.Name, accountDB.Client.Name)
	s.Equal(account.Client.Email, accountDB.Client.Email)
	s.Equal(account.Balance, accountDB.Balance)
}

func (s *AccountDBTestSuite) TestGetWhenAccountDoesNotExist() {
	account, err := s.accountDB.Get("invalid_id")
	s.Error(err)
	s.Nil(account)
}

func TestAccountDBTestSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTestSuite))
}
