package service

import (
	products "github.com/lavatee/shop_products"
	"github.com/lavatee/shop_products/internal/repository"
)

type ProductsService struct {
	Repo *repository.Repository
}

func NewProductsService(repo *repository.Repository) *ProductsService {
	return &ProductsService{
		Repo: repo,
	}
}

func (s *ProductsService) PostProduct(name string, amount int, price int, category string, description string, userId int) (int, error) {
	return s.Repo.PostProduct(name, amount, price, category, description, userId)
}

func (s *ProductsService) DeleteProduct(id int) error {
	return s.Repo.DeleteProduct(id)
}

func (s *ProductsService) GetProducts(category string) ([]products.Product, error) {
	return s.Repo.GetProducts(category)
}

func (s *ProductsService) GetUserProducts(userId int) ([]products.Product, error) {
	return s.Repo.GetUserProducts(userId)
}

func (s *ProductsService) GetSavedProducts(ids []int) ([]products.Product, error) {
	return s.Repo.GetSavedProducts(ids)
}

func (s *ProductsService) GetOneProduct(id int) (products.Product, error) {
	return s.Repo.GetOneProduct(id)
}

func (s *ProductsService) PostOrder(productId int) error {
	return s.Repo.PostOrder(productId)
}
