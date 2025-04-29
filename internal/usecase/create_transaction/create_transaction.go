package createtransaction

import (
	"github.com/guimartiins/eda-go/internal/entity"
	"github.com/guimartiins/eda-go/internal/gateway"
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
}

func NewCreateTransactionUseCase(transactionGateway gateway.TransactionGateway, accountGateway gateway.AccountGateway) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		TransactionGateway: transactionGateway,
		AccountGateway:     accountGateway,
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

	return CreateTransactionOutputDTO{
		ID: transaction.ID,
	}, nil
}
