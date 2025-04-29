package gateway

import "github.com/guimartiins/eda-go/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
