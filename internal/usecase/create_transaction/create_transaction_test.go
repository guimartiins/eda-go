package create_transaction

import (
	"context"
	"testing"

	"github.com/guimartiins/eda-go/internal/entity"
	"github.com/guimartiins/eda-go/internal/event"
	"github.com/guimartiins/eda-go/internal/usecase/mocks"
	"github.com/guimartiins/eda-go/pkg/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TransactionGatewayMock struct {
	mock.Mock
}

func (m *TransactionGatewayMock) Create(transaction *entity.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

type AccountGatewayMock struct {
	mock.Mock
}

func (m *AccountGatewayMock) Save(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *AccountGatewayMock) FindByID(id string) (*entity.Account, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountGatewayMock) UpdateBalance(account *entity.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

type CreateTransactionUseCaseTestSuite struct {
	suite.Suite
	ctx      context.Context
	mockUow  *mocks.UowMock
	useCase  *CreateTransactionUseCase
	account1 *entity.Account
	account2 *entity.Account
}

func (suite *CreateTransactionUseCaseTestSuite) SetupTest() {
	client1, _ := entity.NewClient("client1", "client1@email.com")
	suite.account1 = entity.NewAccount(client1)
	suite.account1.Credit(1000)
	dispatcher := events.NewEventDispatcher()
	eventTransaction := event.NewTransactionCreatedEvent()
	eventBalance := event.NewBalanceUpdatedEvent()

	client2, _ := entity.NewClient("client2", "client2@email.com")
	suite.account2 = entity.NewAccount(client2)
	suite.account2.Credit(1000)

	mockUow := &mocks.UowMock{}

	ctx := context.Background()

	suite.mockUow = mockUow
	suite.ctx = ctx
	suite.useCase = NewCreateTransactionUseCase(mockUow, dispatcher, eventTransaction, eventBalance)
}

func (suite *CreateTransactionUseCaseTestSuite) TestExecute_SuccessfulTransaction() {
	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: suite.account1.ID,
		AccountIDTo:   suite.account2.ID,
		Amount:        100,
	}

	suite.mockUow.On("Do", mock.Anything, mock.Anything).Return(nil)

	output, err := suite.useCase.Execute(suite.ctx, inputDto)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), output)
	suite.mockUow.AssertExpectations(suite.T())
	suite.mockUow.AssertCalled(suite.T(), "Do", mock.Anything, mock.Anything)
}

func (suite *CreateTransactionUseCaseTestSuite) TestExecute_InsufficientFunds() {
	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: suite.account1.ID,
		AccountIDTo:   suite.account2.ID,
		Amount:        2000,
	}
	suite.mockUow.On("Do", mock.Anything, mock.Anything).Return(entity.ErrInsufficientFunds)

	output, err := suite.useCase.Execute(suite.ctx, inputDto)

	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), output)
	assert.Equal(suite.T(), "insufficient funds", err.Error())
}

func TestCreateTransactionUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateTransactionUseCaseTestSuite))
}
