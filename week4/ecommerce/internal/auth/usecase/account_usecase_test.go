package usecase

import (
	"context"
	"ecommerce/internal/auth/entity"
	mock_repository "ecommerce/internal/auth/mocks"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type AccountUsecaseTestSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	mockRepo       *mock_repository.MockIAccountRepository
	accountUsecase IAccountUsecase
}

func (suite *AccountUsecaseTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mock_repository.NewMockIAccountRepository(suite.mockCtrl)
	suite.accountUsecase = NewAccountUsecase(suite.mockRepo)
}

func (suite *AccountUsecaseTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestAccountUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(AccountUsecaseTestSuite))
}

func (suite *AccountUsecaseTestSuite) TestRegister() {
    testCases := []struct {
        name          string
        input         *entity.Account
        mockBehavior  func()
        expectedError error
    }{
        {
            name: "Successful registration",
            input: &entity.Account{
                Username: "testuser",
                Email:    "test@example.com",
                Password: "password123",
            },
            mockBehavior: func() {
                suite.mockRepo.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, nil).Times(2)
                suite.mockRepo.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)
            },
            expectedError: nil,
        },
        {
            name: "Username already exists",
            input: &entity.Account{
                Username: "existinguser",
                Email:    "test@example.com",
                Password: "password123",
            },
            mockBehavior: func() {
                suite.mockRepo.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&entity.Account{}, nil).Times(1)
            },
            expectedError: ErrUsernameExists,
        },
    }

    for _, tc := range testCases {
        suite.Run(tc.name, func() {
            tc.mockBehavior()
            err := suite.accountUsecase.Register(context.Background(), tc.input)
            if tc.expectedError != nil {
                suite.EqualError(err, tc.expectedError.Error())
            } else {
                suite.NoError(err)
            }
        })
    }
}
