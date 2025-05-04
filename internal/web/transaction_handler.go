package web

import (
	"encoding/json"
	"net/http"

	"github.com/guimartiins/eda-go/internal/usecase/create_transaction"
)

type WebTransactionHandler struct {
	CreateTransactionUsecase create_transaction.CreateTransactionUseCase
}

func NewWebTransactionHandler(createTransactionUsecase create_transaction.CreateTransactionUseCase) *WebTransactionHandler {
	return &WebTransactionHandler{
		CreateTransactionUsecase: createTransactionUsecase,
	}
}

func (h *WebTransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var dto create_transaction.CreateTransactionInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	ctx := r.Context()

	output, err := h.CreateTransactionUsecase.Execute(ctx, dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}
	w.WriteHeader(http.StatusCreated)
}
