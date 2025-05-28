package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type GetProductAndStoreIfNeededUseCase struct {
	cacheDuration  time.Duration
	cache          domain.Cache
	serviceRepo    domain.GetProductRepository // source of truth
	databaseRepo   domain.GetProductRepository // local storage
	productStorer  domain.UpsertProductRepository
	productRemover domain.DeleteProductRepository
}

func NewGetProductAndStoreIfNeededUseCase(
	cacheDuration time.Duration,
	cache domain.Cache,
	serviceRepo domain.GetProductRepository,
	databaseRepo domain.GetProductRepository,
	databaseStorer domain.UpsertProductRepository,
	productRemover domain.DeleteProductRepository,
) *GetProductAndStoreIfNeededUseCase {
	return &GetProductAndStoreIfNeededUseCase{
		cacheDuration:  cacheDuration,
		cache:          cache,
		serviceRepo:    serviceRepo,
		databaseRepo:   databaseRepo,
		productStorer:  databaseStorer,
		productRemover: productRemover,
	}
}

func (u *GetProductAndStoreIfNeededUseCase) Execute(ctx context.Context, productID string) (*domain.Product, error) {
	if productID == "" {
		return nil, &e.ValidationError{
			Field: "productID",
			Err:   "is required",
		}
	}

	cacheKey := fmt.Sprintf("product::%s", productID)
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var product domain.Product
		if err := json.Unmarshal([]byte(cached), &product); err == nil {
			return &product, nil
		}
	}

	product, err := u.serviceRepo.GetByID(ctx, productID)

	if err != nil {
		if err == e.NewNotFoundError("product") {
			u.productRemover.Delete(ctx, productID)
		}
		fmt.Printf("[get_product_and_store_if_needed_usecase] ERROR fetching from service: %s", err.Error())
	}

	if product != nil {
		if err := u.productStorer.Upsert(ctx, *product); err != nil {
			fmt.Printf("error upserting product in database: %v\n", err)
		}

		if productJSON, err := json.Marshal(product); err == nil {
			if err := u.cache.Set(ctx, cacheKey, string(productJSON), u.cacheDuration); err != nil {
				fmt.Printf("error storing product in cache: %v\n", err)
			}
		}

		return product, nil
	}

	product, err = u.databaseRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("error fetching from database: %w", err)
	}

	if product != nil {
		if productJSON, err := json.Marshal(product); err == nil {
			if err := u.cache.Set(ctx, cacheKey, string(productJSON), u.cacheDuration); err != nil {
				fmt.Printf("error storing product in cache: %v\n", err)
			}
		}
		return product, nil
	}

	return nil, e.NewNotFoundError(fmt.Sprintf("product %s", productID))
}
