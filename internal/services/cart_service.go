package services

import "go_ecommerce/internal/repositories"

type CartService struct {
	CartRepo *repositories.CartRepository
}

func NewCartService(repo *repositories.CartRepository) *CartService {
	return &CartService{
		CartRepo: repo,
	}
}
