package repository

import (
	"exchangeapp/backend/models"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 通过构造函数注入DB依赖，方便后续单测用 sqlmock 替换
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// Create 创建文章
func (r *ArticleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

// FindByPage 分页查询文章列表
func (r *ArticleRepository) FindByPage(page, size int) (*models.PageResponse, error) {
	var articles []models.Article
	var total int64

	if err := r.db.Model(&models.Article{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := r.db.Limit(size).Offset((page - 1) * size).Find(&articles).Error; err != nil {
		return nil, err
	}

	return &models.PageResponse{
		Data:  articles,
		Page:  page,
		Size:  size,
		Total: total,
	}, nil
}

// FindByID 根据ID查询单篇文章
func (r *ArticleRepository) FindByID(id string) (*models.Article, error) {
	var article models.Article
	if err := r.db.Where("id = ?", id).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}
