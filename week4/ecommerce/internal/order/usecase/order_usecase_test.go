package usecase

import (
	"context"
	"ecommerce/internal/order/entity"
	mock_repository "ecommerce/internal/order/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type OrderUsecaseTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	mockRepo     *mock_repository.MockIOrderRepository
	orderUsecase OrderUsecase
}

func (suite *OrderUsecaseTestSuite) SetupTesT() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mock_repository.NewMockIOrderRepository(suite.mockCtrl)
	suite.orderUsecase = *NewOrderUsecase(suite.mockRepo)
}

func (suite *OrderUsecaseTestSuite) TestDownTest() {
	suite.mockCtrl.Finish()
}

func TestOrderUsecaseSuite(t *testing.T) {
    suite.Run(t, new(OrderUsecaseTestSuite))
}

func (suite *OrderUsecaseTestSuite) TestCreateOrder() {
	testCases := []struct {
		name          string
		input         *entity.Order
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Successful order creation",
			input: &entity.Order{
				UserID: 1,
				Lines: []entity.OrderLine{
					{ProductID: 1, Qty: 2},
					{ProductID: 2, Qty: 1},
				},
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed order creation",
			input: &entity.Order{
				UserID: 2,
				Lines: []entity.OrderLine{
					{ProductID: 3, Qty: 1},
				},
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
			err := suite.orderUsecase.CreateOrder(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetAllOrders() {
	testCases := []struct {
		name           string
		mockBehavior   func()
		expectedResult []*entity.Order
		expectedError  error
	}{
		{
			name: "Successful retrieval of all orders",
			mockBehavior: func() {
				orders := []*entity.Order{
					{ID: 1, UserID: 1, TotalPrice: 100},
					{ID: 2, UserID: 2, TotalPrice: 200},
				}
				suite.mockRepo.EXPECT().GetAll(gomock.Any()).Return(orders, nil)
			},
			expectedResult: []*entity.Order{
				{ID: 1, UserID: 1, TotalPrice: 100},
				{ID: 2, UserID: 2, TotalPrice: 200},
			},
			expectedError: nil,
		},
		{
			name: "Failed retrieval of all orders",
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
			orders, err := suite.orderUsecase.GetAllOrders(context.Background())
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, orders)
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetOrderByID() {
	testCases := []struct {
		name           string
		input          int
		mockBehavior   func()
		expectedResult *entity.Order
		expectedError  error
	}{
		{
			name:  "Successful retrieval of order by ID",
			input: 1,
			mockBehavior: func() {
				order := &entity.Order{ID: 1, UserID: 1, TotalPrice: 100}
				suite.mockRepo.EXPECT().GetByID(gomock.Any(), 1).Return(order, nil)
			},
			expectedResult: &entity.Order{ID: 1, UserID: 1, TotalPrice: 100},
			expectedError:  nil,
		},
		{
			name:  "Failed retrieval of order by ID",
			input: 2,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetByID(gomock.Any(), 2).Return(nil, errors.New("order not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("order not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			order, err := suite.orderUsecase.GetOrderByID(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, order)
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetUserOrders() {
	testCases := []struct {
		name           string
		input          string
		mockBehavior   func()
		expectedResult []*entity.Order
		expectedError  error
	}{
		{
			name:  "Successful retrieval of user orders",
			input: "testuser",
			mockBehavior: func() {
				orders := []*entity.Order{
					{ID: 1, UserID: 1, TotalPrice: 100},
					{ID: 2, UserID: 1, TotalPrice: 200},
				}
				suite.mockRepo.EXPECT().GetUserOrders(gomock.Any(), "testuser").Return(orders, nil)
			},
			expectedResult: []*entity.Order{
				{ID: 1, UserID: 1, TotalPrice: 100},
				{ID: 2, UserID: 1, TotalPrice: 200},
			},
			expectedError: nil,
		},
		{
			name:  "Failed retrieval of user orders",
			input: "nonexistentuser",
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetUserOrders(gomock.Any(), "nonexistentuser").Return(nil, errors.New("user not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("user not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			orders, err := suite.orderUsecase.GetUserOrders(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, orders)
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestUpdateOrder() {
	testCases := []struct {
		name          string
		input         *entity.Order
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Successful order update",
			input: &entity.Order{
				ID:         1,
				UserID:     1,
				TotalPrice: 150,
				OrderDate:  &time.Time{},
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed order update",
			input: &entity.Order{
				ID:         2,
				UserID:     2,
				TotalPrice: 200,
				OrderDate:  &time.Time{},
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("order not found"))
			},
			expectedError: errors.New("order not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			err := suite.orderUsecase.UpdateOrder(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestDeleteOrder() {
	testCases := []struct {
		name          string
		input         int
		mockBehavior  func()
		expectedError error
	}{
		{
			name:  "Successful order deletion",
			input: 1,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "Failed order deletion",
			input: 2,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Delete(gomock.Any(), 2).Return(errors.New("order not found"))
			},
			expectedError: errors.New("order not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			err := suite.orderUsecase.DeleteOrder(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *OrderUsecaseTestSuite) TestGetInvoice() {
	testCases := []struct {
		name           string
		input          int
		mockBehavior   func()
		expectedResult []*entity.InvoiceData
		expectedError  error
	}{
		{
			name:  "Successful invoice retrieval",
			input: 1,
			mockBehavior: func() {
				invoiceData := []*entity.InvoiceData{
					{
						OrderID:      1,
						OrderDate:    "2023-04-14",
						CustomerName: "John Doe",
						Items: []entity.InvoiceItem{
							{ProductName: "Product A", Quantity: 2, UnitPrice: 10, TotalPrice: 20},
						},
						Total: 20,
					},
				}
				suite.mockRepo.EXPECT().GetInvoice(gomock.Any(), 1).Return(invoiceData, nil)
			},
			expectedResult: []*entity.InvoiceData{
				{
					OrderID:      1,
					OrderDate:    "2023-04-14",
					CustomerName: "John Doe",
					Items: []entity.InvoiceItem{
						{ProductName: "Product A", Quantity: 2, UnitPrice: 10, TotalPrice: 20},
					},
					Total: 20,
				},
			},
			expectedError: nil,
		},
		{
			name:  "Failed invoice retrieval",
			input: 2,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetInvoice(gomock.Any(), 2).Return(nil, errors.New("invoice not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("invoice not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			invoice, err := suite.orderUsecase.GetInvoice(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, invoice)
			}
		})
	}
}