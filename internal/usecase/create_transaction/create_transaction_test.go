package create_transaction

import (
	"testing"

	"github.com/guimartiins/eda-go/internal/entity"
	"github.com/guimartiins/eda-go/internal/event"
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

type CreateTransactionUseCaseTestSuite struct {
	suite.Suite
	mockAccount     *AccountGatewayMock
	mockTransaction *TransactionGatewayMock
	useCase         *CreateTransactionUseCase
	account1        *entity.Account
	account2        *entity.Account
}

func (suite *CreateTransactionUseCaseTestSuite) SetupTest() {
	client1, _ := entity.NewClient("client1", "client1@email.com")
	suite.account1 = entity.NewAccount(client1)
	suite.account1.Credit(1000)
	dispatcher := events.NewEventDispatcher()
	event := event.NewTransactionCreatedEvent()

	client2, _ := entity.NewClient("client2", "client2@email.com")
	suite.account2 = entity.NewAccount(client2)
	suite.account2.Credit(1000)

	suite.mockAccount = &AccountGatewayMock{}
	suite.mockTransaction = &TransactionGatewayMock{}
	suite.useCase = NewCreateTransactionUseCase(suite.mockTransaction, suite.mockAccount, dispatcher, event)
}

func (suite *CreateTransactionUseCaseTestSuite) TestExecute_SuccessfulTransaction() {
	suite.mockAccount.On("FindByID", suite.account1.ID).Return(suite.account1, nil)
	suite.mockAccount.On("FindByID", suite.account2.ID).Return(suite.account2, nil)
	suite.mockTransaction.On("Create", mock.Anything).Return(nil)

	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: suite.account1.ID,
		AccountIDTo:   suite.account2.ID,
		Amount:        100,
	}

	output, err := suite.useCase.Execute(inputDto)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), output)
	assert.NotEmpty(suite.T(), output.ID)
	suite.mockAccount.AssertExpectations(suite.T())
	suite.mockTransaction.AssertExpectations(suite.T())
}

func (suite *CreateTransactionUseCaseTestSuite) TestExecute_InsufficientFunds() {
	suite.mockAccount.On("FindByID", suite.account1.ID).Return(suite.account1, nil)
	suite.mockAccount.On("FindByID", suite.account2.ID).Return(suite.account2, nil)

	inputDto := CreateTransactionInputDTO{
		AccountIDFrom: suite.account1.ID,
		AccountIDTo:   suite.account2.ID,
		Amount:        2000,
	}

	output, err := suite.useCase.Execute(inputDto)

	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), entity.ErrInsufficientFunds, err)
	assert.Empty(suite.T(), output.ID)
}

func TestCreateTransactionUseCaseSuite(t *testing.T) {
	suite.Run(t, new(CreateTransactionUseCaseTestSuite))
}
