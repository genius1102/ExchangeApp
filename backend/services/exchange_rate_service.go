package services

import (
	"time"

	"exchangeapp/backend/models"
	"exchangeapp/backend/repository"
)

type ExchangeRateService struct {
	repo *repository.ExchangeRateRepository
}

func NewExchangeRateService(repo *repository.ExchangeRateRepository) *ExchangeRateService {
	return &ExchangeRateService{repo: repo}
}

// CreateExchangeRate 创建汇率记录，自动设置日期
func (s *ExchangeRateService) CreateExchangeRate(rate *models.ExchangeRate) error {
	rate.Date = time.Now()
	return s.repo.Create(rate)
}

// GetExchangeRates 获取全部汇率记录
func (s *ExchangeRateService) GetExchangeRates() ([]models.ExchangeRate, error) {
	return s.repo.FindAll()
}
