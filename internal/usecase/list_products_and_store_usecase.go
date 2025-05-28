package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ydoro/wishlist/internal/domain"
	e "github.com/ydoro/wishlist/internal/domain/errors"
)

type ListProductsAndStoreUseCase struct {
	cacheDuration  time.Duration
	cache          domain.Cache
	serviceRepo    domain.ListProductsRepository // source of truth
	databaseRepo   domain.ListProductsRepository // local storage
	productStorer  domain.UpsertProductRepository
	productRemover domain.DeleteProductRepository
}

func NewListProductsAndStoreUseCase(
	cacheDuration time.Duration,
	cache domain.Cache,
	serviceRepo domain.ListProductsRepository,
	databaseRepo domain.ListProductsRepository,
	databaseStorer domain.UpsertProductRepository,
	productRemover domain.DeleteProductRepository,
) *ListProductsAndStoreUseCase {
	return &ListProductsAndStoreUseCase{
		cacheDuration:  cacheDuration,
		cache:          cache,
		serviceRepo:    serviceRepo,
		databaseRepo:   databaseRepo,
		productStorer:  databaseStorer,
		productRemover: productRemover,
	}
}

func (u *ListProductsAndStoreUseCase) Execute(ctx context.Context, count int, offset int) (*[]domain.Product, error) {
	cacheKey := fmt.Sprintf("products::%d::%d", count, offset)
	if cached, err := u.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var products []domain.Product
		if err := json.Unmarshal([]byte(cached), &products); err == nil {
			return &products, nil
		}
	}

	products, err := u.serviceRepo.List(ctx, count, offset)

	if err != nil {
		fmt.Printf("[list_products_and_storeusecase] ERROR fetching from service: %s", err.Error())
	}

	if products != nil {
		u.upsertProductsFailsafe(ctx, &products)
		if productJSON, err := json.Marshal(products); err == nil {
			if err := u.cache.Set(ctx, cacheKey, string(productJSON), u.cacheDuration); err != nil {
				fmt.Printf("error storing products in cache: %v\n", err)
			}
		}
		return &products, nil
	}

	products, err = u.databaseRepo.List(ctx, count, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching from database: %w", err)
	}

	if products != nil {
		if productJSON, err := json.Marshal(products); err == nil {
			if err := u.cache.Set(ctx, cacheKey, string(productJSON), u.cacheDuration); err != nil {
				fmt.Printf("error storing product in cache: %v\n", err)
			}
		}
		return &products, nil
	}

	return nil, e.NewNotFoundError("products")
}

func (u *ListProductsAndStoreUseCase) upsertProductsFailsafe(ctx context.Context, products *[]domain.Product) {
	var wg sync.WaitGroup
	for _, product := range *products {
		wg.Add(1)
		go func(product domain.Product) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				err := u.productStorer.Upsert(ctx, product)
				if err != nil {
					fmt.Printf("error upserting product in database: %v\n", err)
				}
			}
		}(product)
	}
	wg.Wait()
}
