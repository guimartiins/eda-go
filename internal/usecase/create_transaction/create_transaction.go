package create_transaction

import (
	"github.com/guimartiins/eda-go/internal/entity"
	"github.com/guimartiins/eda-go/internal/gateway"
	"github.com/guimartiins/eda-go/pkg/events"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string
	AccountIDTo   string
	Amount        float64
}

type CreateTransactionOutputDTO struct {
	ID string
}

type CreateTransactionUseCase struct {
	TransactionGateway gateway.TransactionGateway
	AccountGateway     gateway.AccountGateway
	EventDispatcher    events.EventDispatcherInterface
	transactionCreated events.EventInterface
}

func NewCreateTransactionUseCase(transactionGateway gateway.TransactionGateway, accountGateway gateway.AccountGateway, eventDispatcher events.EventDispatcherInterface, transactionCreated events.EventInterface) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: transactionGateway,
		AccountGateway:     accountGateway,
		EventDispatcher:    eventDispatcher,
		transactionCreated: transactionCreated,
	}
}

func (uc *CreateTransactionUseCase) Execute(input CreateTransactionInputDTO) (CreateTransactionOutputDTO, error) {
	accountFrom, err := uc.AccountGateway.FindByID(input.AccountIDFrom)
	if err != nil {
		return CreateTransactionOutputDTO{}, err
	}

	accountTo, err := uc.AccountGateway.FindByID(input.AccountIDTo)
	if err != nil {
		return CreateTransactionOutputDTO{}, err
	}

	transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
	if err != nil {
		return CreateTransactionOutputDTO{}, err
	}

	err = uc.TransactionGateway.Create(transaction)
	if err != nil {
		return CreateTransactionOutputDTO{}, err
	}

	output := CreateTransactionOutputDTO{
		ID: transaction.ID,
	}

	uc.transactionCreated.SetPayload(output)
	uc.EventDispatcher.Dispatch(uc.transactionCreated)

	return output, nil
}
