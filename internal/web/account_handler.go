package web

import (
	"encoding/json"
	"net/http"

	"github.com/guimartiins/eda-go/internal/usecase/create_account"
)

type WebAccountHandler struct {
	CreateAccountUsecase create_account.CreateAccountUseCase
}

func NewWebAccountHandler(createAccountUsecase create_account.CreateAccountUseCase) *WebAccountHandler {
	return &WebAccountHandler{
		CreateAccountUsecase: createAccountUsecase,
	}
}

func (h *WebAccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var dto create_account.CreateAccountInputDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.CreateAccountUsecase.Execute(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
