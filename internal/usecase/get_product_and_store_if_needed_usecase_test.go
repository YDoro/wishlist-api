package usecase_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
	"github.com/ydoro/wishlist/internal/usecase"
	mocks "github.com/ydoro/wishlist/mock/domain"
	"go.uber.org/mock/gomock"
)

func TestGetProductAndStoreIfNeededUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCache(ctrl)
	serviceRepo := mocks.NewMockGetProductRepository(ctrl)
	databaseRepo := mocks.NewMockGetProductRepository(ctrl)
	databaseStorer := mocks.NewMockUpsertProductRepository(ctrl)
	productRemover := mocks.NewMockDeleteProductRepository(ctrl)
	mockedTime := time.Now()
	tests := []struct {
		name          string
		setupMocks    func()
		productId     string
		expectError   error
		expectProduct *domain.Product
	}{
		{
			name:       "should return error if product id is empty",
			setupMocks: func() {},
			productId:  "",
			expectError: &e.ValidationError{
				Field: "productID",
				Err:   "is required",
			},
		},
		{
			name: "should a product from cache",
			setupMocks: func() {
				productJSON, _ := json.Marshal(domain.Product{
					ID: "123",
				})
				mockCache.EXPECT().Get(gomock.Any(), "product::123").Return(string(productJSON), nil)
			},
			productId: "123",
			expectProduct: &domain.Product{
				ID: "123",
			},
		},
		{
			name: "should store a product in cache and database if service returns a product",
			setupMocks: func() {
				p := &domain.Product{
					ID: "123",
				}
				mockCache.EXPECT().Get(gomock.Any(), "product::123").Return("", nil)
				serviceRepo.EXPECT().GetByID(gomock.Any(), "123").Return(p, nil)
				databaseStorer.EXPECT().Upsert(gomock.Any(), *p).Return(nil)
				mockCache.EXPECT().Set(gomock.Any(), "product::123", gomock.Any(), gomock.Any()).Return(nil)
			},
			productId: "123",
			expectProduct: &domain.Product{
				ID: "123",
			},
		},
		{
			name: "should try to delete a product on product not found error",
			setupMocks: func() {
				mockCache.EXPECT().Get(gomock.Any(), "product::123").Return("", nil)
				serviceRepo.EXPECT().GetByID(gomock.Any(), "123").Return(nil, e.NewNotFoundError("product"))
				productRemover.EXPECT().Delete(gomock.Any(), "123").Return(nil)
				databaseRepo.EXPECT().GetByID(gomock.Any(), "123").Return(&domain.Product{
					ID:        "123",
					DeletedAt: mockedTime.Format(time.RFC3339),
				}, nil)
				mockCache.EXPECT().Set(gomock.Any(), "product::123", gomock.Any(), gomock.Any()).Return(fmt.Errorf("error storing product in cache"))
			},
			productId: "123",
			expectProduct: &domain.Product{
				ID:        "123",
				DeletedAt: mockedTime.Format(time.RFC3339),
			},
		},
		{
			name: "should return the product even if the upsert or cache fails",
			setupMocks: func() {
				p := &domain.Product{
					ID: "123",
				}
				mockCache.EXPECT().Get(gomock.Any(), "product::123").Return("", nil)
				serviceRepo.EXPECT().GetByID(gomock.Any(), "123").Return(p, nil)
				databaseStorer.EXPECT().Upsert(gomock.Any(), *p).Return(fmt.Errorf("error upserting product in database"))
				mockCache.EXPECT().Set(gomock.Any(), "product::123", gomock.Any(), gomock.Any()).Return(fmt.Errorf("error storing product in cache"))
			},
			productId: "123",
			expectProduct: &domain.Product{
				ID: "123",
			},
		},
		{
			name: "should return error if cache service and database fails",
			setupMocks: func() {
				mockCache.EXPECT().Get(gomock.Any(), "product::123").Return("", nil)
				serviceRepo.EXPECT().GetByID(gomock.Any(), "123").Return(nil, nil)
				databaseRepo.EXPECT().GetByID(gomock.Any(), "123").Return(nil, fmt.Errorf("foo"))
			},
			productId:   "123",
			expectError: fmt.Errorf("error fetching from database: %w", fmt.Errorf("foo")),
		},
		{
			name: "should return product not found if everything fails with nil nil",
			setupMocks: func() {
				mockCache.EXPECT().Get(gomock.Any(), "product::123").Return("", nil)
				serviceRepo.EXPECT().GetByID(gomock.Any(), "123").Return(nil, nil)
				databaseRepo.EXPECT().GetByID(gomock.Any(), "123").Return(nil, nil)
			},
			productId:   "123",
			expectError: e.NewNotFoundError("product 123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			sut := usecase.NewGetProductAndStoreIfNeededUseCase(time.Second, mockCache, serviceRepo, databaseRepo, databaseStorer, productRemover)
			p, err := sut.Execute(context.Background(), tt.productId)

			assert.Equal(t, tt.expectError, err)
			assert.Equal(t, tt.expectProduct, p)
		})
	}
}
