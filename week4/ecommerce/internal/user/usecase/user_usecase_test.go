package usecase

import (
	"context"
	"ecommerce/internal/user/entity"
	mock_repository "ecommerce/internal/user/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserUsecaseTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockRepo    *mock_repository.MockIUser
	userUsecase UserUsecase
}

func (suite *UserUsecaseTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mock_repository.NewMockIUser(suite.mockCtrl)
	suite.userUsecase = *NewUserUsecase(suite.mockRepo)
}

func (suite *UserUsecaseTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestUserUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUsecaseTestSuite))
}

func (suite *UserUsecaseTestSuite) TestCreateUser() {
	testCases := []struct {
		name          string
		input         *entity.User
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Successful user creation",
			input: &entity.User{
				Name:     "John Doe",
				Username: "johndoe",
				Email:    "john@example.com",
				Balance:  100.0,
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed user creation",
			input: &entity.User{
				Name:     "Jane Doe",
				Username: "janedoe",
				Email:    "jane@example.com",
				Balance:  50.0,
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			err := suite.userUsecase.CreateUser(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetAllUsers() {
	testCases := []struct {
		name           string
		mockBehavior   func()
		expectedResult []*entity.User
		expectedError  error
	}{
		{
			name: "Successful retrieval of users",
			mockBehavior: func() {
				users := []*entity.User{
					{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0},
					{ID: 2, Name: "Jane Doe", Username: "janedoe", Email: "jane@example.com", Balance: 50.0},
				}
				suite.mockRepo.EXPECT().GetAll(gomock.Any()).Return(users, nil)
			},
			expectedResult: []*entity.User{
				{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0},
				{ID: 2, Name: "Jane Doe", Username: "janedoe", Email: "jane@example.com", Balance: 50.0},
			},
			expectedError: nil,
		},
		{
			name: "Failed retrieval of users",
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("database error"))
			},
			expectedResult: nil,
			expectedError:  errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			users, err := suite.userUsecase.GetAllUsers(context.Background())
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, users)
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetByUserID() {
	testCases := []struct {
		name           string
		input          int
		mockBehavior   func()
		expectedResult *entity.User
		expectedError  error
	}{
		{
			name:  "Successful retrieval of user by ID",
			input: 1,
			mockBehavior: func() {
				user := &entity.User{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0}
				suite.mockRepo.EXPECT().GetByID(gomock.Any(), 1).Return(user, nil)
			},
			expectedResult: &entity.User{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0},
			expectedError:  nil,
		},
		{
			name:  "Failed retrieval of user by ID",
			input: 2,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetByID(gomock.Any(), 2).Return(nil, errors.New("user not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			user, err := suite.userUsecase.GetByUserID(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, user)
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetByUserUsername() {
	testCases := []struct {
		name           string
		input          string
		mockBehavior   func()
		expectedResult *entity.User
		expectedError  error
	}{
		{
			name:  "Successful retrieval of user by username",
			input: "johndoe",
			mockBehavior: func() {
				user := &entity.User{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0}
				suite.mockRepo.EXPECT().GetByUsername(gomock.Any(), "johndoe").Return(user, nil)
			},
			expectedResult: &entity.User{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0},
			expectedError:  nil,
		},
		{
			name:  "Failed retrieval of user by username",
			input: "nonexistent",
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetByUsername(gomock.Any(), "nonexistent").Return(nil, errors.New("user not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			user, err := suite.userUsecase.GetByUserUsername(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, user)
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestGetByUserEmail() {
	testCases := []struct {
		name           string
		input          string
		mockBehavior   func()
		expectedResult *entity.User
		expectedError  error
	}{
		{
			name:  "Successful retrieval of user by email",
			input: "john@example.com",
			mockBehavior: func() {
				user := &entity.User{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0}
				suite.mockRepo.EXPECT().GetByEmail(gomock.Any(), "john@example.com").Return(user, nil)
			},
			expectedResult: &entity.User{ID: 1, Name: "John Doe", Username: "johndoe", Email: "john@example.com", Balance: 100.0},
			expectedError:  nil,
		},
		{
			name:  "Failed retrieval of user by email",
			input: "nonexistent@example.com",
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetByEmail(gomock.Any(), "nonexistent@example.com").Return(nil, errors.New("user not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			user, err := suite.userUsecase.GetByUserEmail(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, user)
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestUpdateUser() {
	testCases := []struct {
		name          string
		input         *entity.User
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Successful user update",
			input: &entity.User{
				ID:       1,
				Name:     "John Doe Updated",
				Username: "johndoe_updated",
				Email:    "john_updated@example.com",
				Balance:  150.0,
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed user update",
			input: &entity.User{
				ID:       2,
				Name:     "Jane Doe Updated",
				Username: "janedoe_updated",
				Email:    "jane_updated@example.com",
				Balance:  75.0,
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			err := suite.userUsecase.UpdateUser(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *UserUsecaseTestSuite) TestDeleteUser() {
	testCases := []struct {
		name          string
		input         int
		mockBehavior  func()
		expectedError error
	}{
		{
			name:  "Successful user deletion",
			input: 1,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "Failed user deletion",
			input: 2,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Delete(gomock.Any(), 2).Return(errors.New("user not found"))
			},
			expectedError: errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			err := suite.userUsecase.DeleteUser(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}
