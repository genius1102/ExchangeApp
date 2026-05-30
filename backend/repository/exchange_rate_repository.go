package repository

import (
	"exchangeapp/backend/models"

	"gorm.io/gorm"
)

type ExchangeRateRepository struct {
	db *gorm.DB
}

func NewExchangeRateRepository(db *gorm.DB) *ExchangeRateRepository {
	return &ExchangeRateRepository{db: db}
}

func (r *ExchangeRateRepository) Create(rate *models.ExchangeRate) error {
	return r.db.Create(rate).Error
}

func (r *ExchangeRateRepository) FindAll() ([]models.ExchangeRate, error) {
	var rates []models.ExchangeRate
	if err := r.db.Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}
