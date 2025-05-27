package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type GetProductAndStoreIfNeededUseCase struct {
	cacheDuration time.Duration
	cache         domain.Cache
	FirstGetter   domain.GetProductRepository
	SecondGetter  domain.GetProductRepository
	Storer        domain.CreateProductRepository
}

func NewGetProductAndStoreIfNeededUseCase(
	cacheDuration time.Duration,
	cacheGetter domain.Cache,
	firstGetter domain.GetProductRepository,
	secondGetter domain.GetProductRepository,
	storer domain.CreateProductRepository) *GetProductAndStoreIfNeededUseCase {
	return &GetProductAndStoreIfNeededUseCase{
		cache:        cacheGetter,
		FirstGetter:  firstGetter,
		SecondGetter: secondGetter,
		Storer:       storer,
	}
}

func (u *GetProductAndStoreIfNeededUseCase) Execute(ctx context.Context, productID string) (*domain.Product, error) {
	var internalErr error
	storeErr := func(err error) {
		if err != nil {
			errors.Join(internalErr, err)
		}
	}

	if productID == "" {
		return nil, &e.ValidationError{
			Field: "productID",
			Err:   "is required",
		}
	}

	product := &domain.Product{}
	cached, err := u.cache.Get(ctx, fmt.Sprintf("%s::%s", "product", productID))
	storeErr(err)

	err = json.Unmarshal([]byte(cached), product)
	storeErr(err)

	if product.ID != "" {
		return product, nil
	}

	product, err = u.FirstGetter.GetByID(ctx, productID)
	storeErr(err)

	if product != nil {
		pAsJson, err := json.Marshal(product)
		if err == nil {
			u.cache.Set(ctx, fmt.Sprintf("%s::%s", "product", productID), string(pAsJson), u.cacheDuration)
		}
		return product, nil
	}

	product, err = u.SecondGetter.GetByID(ctx, productID)

	if err != nil {
		storeErr(err)
		return nil, internalErr
	}

	if product != nil {
		pAsJson, err := json.Marshal(product)
		if err == nil {
			u.cache.Set(ctx, fmt.Sprintf("%s::%s", "product", productID), string(pAsJson), u.cacheDuration)
		}
		return product, nil
	}

	return nil, internalErr
}
