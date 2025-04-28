package gateway

import "github.com/guimartiins/fcutils/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
