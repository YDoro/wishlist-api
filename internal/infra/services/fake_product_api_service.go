package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ydoro/wishlist/internal/domain"
)

type FakeProductAPIService struct {
	baseurl string
	client  domain.HttpClient
}

func NewFakeProductAPIService(baseurl string, client domain.HttpClient) *FakeProductAPIService {
	return &FakeProductAPIService{
		baseurl: baseurl,
		client:  client,
	}
}

type fakeProduct struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
	Rating      struct {
		Rate  float64 `json:"rate"`
		Count int     `json:"count"`
	} `json:"rating"`
}

func fakeProductToDomain(p fakeProduct) domain.Product {
	return domain.Product{
		ID:          strconv.Itoa(p.ID),
		Name:        p.Title,
		Description: p.Description,
		Price:       float64(p.Price),
		Category:    p.Category,
		Images:      []string{p.Image},
		Rating: &domain.Rating{
			Average: p.Rating.Rate,
			Count:   p.Rating.Count,
		},
	}
}

func (ps *FakeProductAPIService) GetByID(ctx context.Context, productID string) (*domain.Product, error) {
	url := fmt.Sprintf("%s/%s", ps.baseurl, productID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := ps.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	var p fakeProduct
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}

	res := fakeProductToDomain(p)

	return &res, nil
}

func (ps *FakeProductAPIService) List(ctx context.Context, count int, offset int) ([]domain.Product, error) {
	url := fmt.Sprintf("%s?count=%d&offset=%d", ps.baseurl, count, offset)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	resp, err := ps.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var products []*fakeProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	res := make([]domain.Product, len(products))

	for i, p := range products {
		res[i] = fakeProductToDomain(*p)
	}

	return res, nil
}
