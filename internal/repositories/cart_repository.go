package repositories

import (
	"go_ecommerce/database"
	"go_ecommerce/internal/models"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository() *CartRepository {
	return &CartRepository{
		db: database.DB,
	}
}

func (r *CartRepository) AddToCart(cart *models.Cart) error {
	return r.db.Create(cart).Error
}
func (r *CartRepository) RemoveFromCart(userID, productID uint) error {
	return r.db.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Cart{}).Error
}
func (r *CartRepository) GetCartByUserID(userID uint) ([]models.Cart, error) {
	var cart []models.Cart
	err := r.db.Where("user_id = ?", userID).Find(&cart).Error
	return cart, err
}
func (cr *CartRepository) FindByUserAndProduct(userID, productID uint) (*models.Cart, error) {
	var cart models.Cart
	err := cr.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}
func (cr *CartRepository) Create(cart *models.Cart) error {
	return cr.db.Create(cart).Error
}
func (cr *CartRepository) Update(cart *models.Cart) error {
	return cr.db.Save(cart).Error
}
func (cr *CartRepository) Delete(cart *models.Cart) error {
	return cr.db.Delete(cart).Error
}
