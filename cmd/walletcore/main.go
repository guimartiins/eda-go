package main

import (
	"context"
	"database/sql"
	"fmt"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/guimartiins/eda-go/internal/database"
	"github.com/guimartiins/eda-go/internal/event"
	"github.com/guimartiins/eda-go/internal/event/handler"
	create_account "github.com/guimartiins/eda-go/internal/usecase/create_account"
	create_client "github.com/guimartiins/eda-go/internal/usecase/create_client"
	create_transaction "github.com/guimartiins/eda-go/internal/usecase/create_transaction"
	"github.com/guimartiins/eda-go/internal/web"
	"github.com/guimartiins/eda-go/internal/web/webserver"
	"github.com/guimartiins/eda-go/pkg/events"
	"github.com/guimartiins/eda-go/pkg/kafka"
	"github.com/guimartiins/eda-go/pkg/uow"
)

func main() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"root",      // username
		"root",      // password
		"mysql", // host (changed from "mysql" to "localhost")
		"3306",      // port
		"wallet"))   // database name

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("error connecting to the database: %v", err))
	}
	fmt.Println("Successfully connected to database")

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	transactionCreatedEvent := event.NewTransactionCreatedEvent()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	},
	)
	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	},
	)

	createClientUseCase := create_client.NewCreateClientUseCase(clientDb)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDb, clientDb)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent)

	webserver := webserver.NewWebServer("8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)

	fmt.Println("Starting web server")
	webserver.Start()
}
