package database

import (
	"database/sql"

	"github.com/guimartiins/eda-go/internal/entity"
)

type AccountDB struct {
	db *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{db: db}
}

func (a *AccountDB) Get(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client

	stmt, err := a.db.Prepare("SELECT a.id, a.client_id, a.balance, a.created_at, c.id, c.name, c.email, c.created_at FROM accounts a INNER JOIN clients c ON a.client_id = c.id WHERE a.id = ?")

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.ID,
		&account.Client.ID,
		&account.Balance,
		&account.CreatedAt,
		&account.Client.ID,
		&account.Client.Name,
		&account.Client.Email,
		&account.Client.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AccountDB) Save(account *entity.Account) error {
	stmt, err := a.db.Prepare("INSERT INTO accounts (id, client_id, balance, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(account.ID, account.Client.ID, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}
