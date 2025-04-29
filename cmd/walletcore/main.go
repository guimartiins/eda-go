package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guimartiins/eda-go/internal/database"
	"github.com/guimartiins/eda-go/internal/event"
	createaccount "github.com/guimartiins/eda-go/internal/usecase/create_account"
	createclient "github.com/guimartiins/eda-go/internal/usecase/create_client"
	createtransaction "github.com/guimartiins/eda-go/internal/usecase/create_transaction"
	"github.com/guimartiins/eda-go/pkg/events"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", "root", "root", "mysql", "3306", "wallet"))

	if err != nil {
		panic(err)
	}
	defer db.Close()

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreatedEvent()
	// eventDispatcher.Register("TransactionCreated", handler)

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)
	transactionDb := database.NewTransactionDB(db)

	createClientUseCase := createclient.NewCreateClientUseCase(clientDb)
	createAccountUseCase := createaccount.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := createtransaction.NewCreateTransactionUseCase(transactionDb, accountDb, eventDispatcher, transactionCreatedEvent)
}
