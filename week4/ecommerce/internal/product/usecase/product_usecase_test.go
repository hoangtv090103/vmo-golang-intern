package usecase

import (
	"context"
	"ecommerce/internal/product/entity"
	mock_repository "ecommerce/internal/product/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ProductUsecaseTestSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	mockRepo       *mock_repository.MockIProductRepository
	productUsecase ProductUsecase
}

func (suite *ProductUsecaseTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mock_repository.NewMockIProductRepository(suite.mockCtrl)
	suite.productUsecase = *NewProductUsecase(suite.mockRepo)
}

func (suite *ProductUsecaseTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestProductUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUsecaseTestSuite))
}

func (suite *ProductUsecaseTestSuite) TestCreateProduct() {
    testCases := []struct {
        name          string
        input         *entity.Product
        mockBehavior  func()
        expectedError error
    }{
        {
            name: "Successful product creation",
            input: entity.NewProduct().
                SetName("Product A").
                SetDescription("Product description").
                SetStock(100).
                SetPrice(10.0),
            mockBehavior: func() {
                suite.mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
            },
            expectedError: nil,
        },
        {
            name: "Failed product creation - Invalid price",
            input: entity.NewProduct().
                SetName("Product C").
                SetDescription("Product with invalid price").
                SetStock(75).
                SetPrice(-5.0),
            mockBehavior: func() {
                // No mock expectation, as Create should not be called
            },
            expectedError: errors.New("invalid price"),
        },
        {
            name: "Failed product creation - Invalid stock",
            input: entity.NewProduct().
                SetName("Product D").
                SetDescription("Product with invalid stock").
                SetStock(-75).
                SetPrice(5.0),
            mockBehavior: func() {
                // No mock expectation, as Create should not be called
            },
            expectedError: errors.New("invalid stock"),
        },
    }

    for _, tc := range testCases {
        suite.Run(tc.name, func() {
            tc.mockBehavior()
            err := suite.productUsecase.CreateProduct(context.Background(), tc.input)
            if tc.expectedError != nil {
                suite.EqualError(err, tc.expectedError.Error(), "The error does not match the expected error")
            } else {
                suite.NoError(err)
            }
        })
    }
}

func (suite *ProductUsecaseTestSuite) TestGetAllProducts() {
	testCases := []struct {
		name           string
		mockBehavior   func()
		expectedResult []*entity.Product
		expectedError  error
	}{
		{
			name: "Successful retrieval of all products",
			mockBehavior: func() {
				products := []*entity.Product{
					{ID: 1, Name: "Product A", Price: 10.0, Stock: 100},
					{ID: 2, Name: "Product B", Price: 20.0, Stock: 50},
				}
				suite.mockRepo.EXPECT().GetAll(gomock.Any()).Return(products, nil)
			},
			expectedResult: []*entity.Product{
				{ID: 1, Name: "Product A", Price: 10.0, Stock: 100},
				{ID: 2, Name: "Product B", Price: 20.0, Stock: 50},
			},
			expectedError: nil,
		},
		{
			name: "Failed retrieval of all products",
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
			products, err := suite.productUsecase.GetAllProducts(context.Background())
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, products)
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestGetByProductID() {
	testCases := []struct {
		name           string
		input          int
		mockBehavior   func()
		expectedResult *entity.Product
		expectedError  error
	}{
		{
			name:  "Successful retrieval of product by ID",
			input: 1,
			mockBehavior: func() {
				product := &entity.Product{ID: 1, Name: "Product A", Price: 10.0, Stock: 100}
				suite.mockRepo.EXPECT().GetByID(gomock.Any(), 1).Return(product, nil)
			},
			expectedResult: &entity.Product{ID: 1, Name: "Product A", Price: 10.0, Stock: 100},
			expectedError:  nil,
		},
		{
			name:  "Failed retrieval of product by ID",
			input: 2,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetByID(gomock.Any(), 2).Return(nil, errors.New("product not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("product not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			product, err := suite.productUsecase.GetByProductID(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, product)
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestGetByProductName() {
	testCases := []struct {
		name           string
		input          string
		mockBehavior   func()
		expectedResult []*entity.Product
		expectedError  error
	}{
		{
			name:  "Successful retrieval of products by name",
			input: "Product A",
			mockBehavior: func() {
				products := []*entity.Product{{ID: 1, Name: "Product A", Price: 10.0, Stock: 100}}
				suite.mockRepo.EXPECT().GetByName(gomock.Any(), "Product A").Return(products, nil)
			},
			expectedResult: []*entity.Product{{ID: 1, Name: "Product A", Price: 10.0, Stock: 100}},
			expectedError:  nil,
		},
		{
			name:  "Failed retrieval of products by name",
			input: "Nonexistent Product",
			mockBehavior: func() {
				suite.mockRepo.EXPECT().GetByName(gomock.Any(), "Nonexistent Product").Return(nil, errors.New("products not found"))
			},
			expectedResult: nil,
			expectedError:  errors.New("products not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			products, err := suite.productUsecase.GetByProductName(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
				suite.Equal(tc.expectedResult, products)
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestUpdateProduct() {
	testCases := []struct {
		name          string
		input         *entity.Product
		mockBehavior  func()
		expectedError error
	}{
		{
			name: "Successful product update",
			input: &entity.Product{
				ID:    1,
				Name:  "Updated Product A",
				Price: 15.0,
				Stock: 150,
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed product update",
			input: &entity.Product{
				ID:    2,
				Name:  "Updated Product B",
				Price: 25.0,
				Stock: 75,
			},
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errors.New("product not found"))
			},
			expectedError: errors.New("product not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			err := suite.productUsecase.UpdateProduct(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *ProductUsecaseTestSuite) TestDeleteProduct() {
	testCases := []struct {
		name          string
		input         int
		mockBehavior  func()
		expectedError error
	}{
		{
			name:  "Successful product deletion",
			input: 1,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "Failed product deletion",
			input: 2,
			mockBehavior: func() {
				suite.mockRepo.EXPECT().Delete(gomock.Any(), 2).Return(errors.New("product not found"))
			},
			expectedError: errors.New("product not found"),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.mockBehavior()
			err := suite.productUsecase.DeleteProduct(context.Background(), tc.input)
			if tc.expectedError != nil {
				suite.EqualError(err, tc.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}
