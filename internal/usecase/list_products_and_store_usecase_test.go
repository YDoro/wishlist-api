package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
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

func TestListProductsAndStoreUsecase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cacheMock := mocks.NewMockCache(ctrl)
	serviceRepoMock := mocks.NewMockListProductsRepository(ctrl)
	databaseRepoMock := mocks.NewMockListProductsRepository(ctrl)
	productStoreMock := mocks.NewMockUpsertProductRepository(ctrl)
	productRemoverMock := mocks.NewMockDeleteProductRepository(ctrl)
	l1 := []domain.Product{
		{
			ID:   "1",
			Name: "Product 1",
		},
		{
			ID:   "2",
			Name: "Product 2",
		},
	}

	tests := []struct {
		name          string
		mockSetup     func()
		offset        int
		count         int
		expectedError error
		expectedLen   int
	}{
		{
			name:   "should return a list of products from cache",
			offset: 0,
			count:  10,
			mockSetup: func() {
				jL, _ := json.Marshal(l1)
				cacheMock.EXPECT().Get(gomock.Any(), gomock.Any()).Return(string(jL), nil)
			},
			expectedError: nil,
			expectedLen:   2,
		},
		{
			name:   "should fetch products from service and store in cache",
			offset: 0,
			count:  10,
			mockSetup: func() {
				jL, _ := json.Marshal(l1)
				cacheMock.EXPECT().Get(gomock.Any(), "products::10::0").Return("", errors.New("cache error"))
				serviceRepoMock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(l1, nil)
				productStoreMock.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(nil)
				productStoreMock.EXPECT().Upsert(gomock.Any(), gomock.Any()).Return(errors.New("store error"))
				cacheMock.EXPECT().Set(gomock.Any(), "products::10::0", string(jL), gomock.Any()).Return(errors.New("cache error"))

			},
			expectedError: nil,
			expectedLen:   2,
		},
		{
			name:   "should fetch products from database if service fails",
			offset: 0,
			count:  10,
			mockSetup: func() {
				jL, _ := json.Marshal(l1)
				cacheMock.EXPECT().Get(gomock.Any(), "products::10::0").Return("", errors.New("cache error"))
				serviceRepoMock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
				databaseRepoMock.EXPECT().List(gomock.Any(), 10, 0).Return(l1, nil)
				cacheMock.EXPECT().Set(gomock.Any(), "products::10::0", string(jL), gomock.Any()).Return(errors.New("cache error"))
			},
			expectedError: nil,
			expectedLen:   2,
		},
		{
			name:   "should fail if database fails",
			offset: 0,
			count:  10,
			mockSetup: func() {
				cacheMock.EXPECT().Get(gomock.Any(), "products::10::0").Return("", errors.New("cache error"))
				serviceRepoMock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
				databaseRepoMock.EXPECT().List(gomock.Any(), 10, 0).Return(nil, errors.New("database error"))
			},
			expectedError: fmt.Errorf("error fetching from database: %w", errors.New("database error")),
			expectedLen:   0,
		},
		{
			name:   "should fail if database return nil nil",
			offset: 0,
			count:  10,
			mockSetup: func() {
				cacheMock.EXPECT().Get(gomock.Any(), "products::10::0").Return("", errors.New("cache error"))
				serviceRepoMock.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
				databaseRepoMock.EXPECT().List(gomock.Any(), 10, 0).Return(nil, nil)
			},
			expectedError: e.NewNotFoundError("products"),
			expectedLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			sut := usecase.NewListProductsAndStoreUseCase(time.Second, cacheMock, serviceRepoMock, databaseRepoMock, productStoreMock, productRemoverMock)
			ps, err := sut.Execute(context.Background(), tt.count, tt.offset)
			assert.Equal(t, tt.expectedError, err)

			if ps != nil {
				assert.Equal(t, tt.expectedLen, len(*ps))
			} else {
				assert.Equal(t, tt.expectedLen, 0)
			}

		})

	}
}
