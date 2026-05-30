package services

import (
	"encoding/json"
	"fmt"
	"time"

	"exchangeapp/backend/models"
	"exchangeapp/backend/repository"

	"github.com/go-redis/redis"
)

const (
	defaultPageSize = 10
	maxCachePages   = 3
	cacheTTL        = 10 * time.Minute
)

type ArticleService struct {
	repo  *repository.ArticleRepository
	redis *redis.Client
}

// NewArticleService 通过构造函数注入依赖，Controller 只依赖 Service，不感知 DB/Redis
func NewArticleService(repo *repository.ArticleRepository, redis *redis.Client) *ArticleService {
	return &ArticleService{repo: repo, redis: redis}
}

// CreateArticle 创建文章并失效前N页缓存
func (s *ArticleService) CreateArticle(article *models.Article) error {
	if err := s.repo.Create(article); err != nil {
		return err
	}

	// 失效缓存：新文章出现后，前N页的缓存数据已不准确
	// 失败仅打日志，不影响主流程（缓存有TTL兜底）
	for i := 0; i < maxCachePages; i++ {
		cacheKey := s.pageCacheKey(i+1, defaultPageSize)
		if err := s.redis.Del(cacheKey).Err(); err != nil {
			fmt.Printf("warn: failed to delete cache key %s: %v\n", cacheKey, err)
		}
	}

	return nil
}

// GetArticles 分页获取文章列表，前N页走Cache-Aside，后续直接查库
func (s *ArticleService) GetArticles(page, size int) (*models.PageResponse, error) {
	if page <= maxCachePages {
		return s.getArticlesWithCache(page, size)
	}
	return s.repo.FindByPage(page, size)
}

// GetArticleByID 根据ID获取单篇文章
func (s *ArticleService) GetArticleByID(id string) (*models.Article, error) {
	return s.repo.FindByID(id)
}

// ========== 私有方法 ==========

// getArticlesWithCache Cache-Aside模式：先查Redis，未命中则查DB并回填缓存
func (s *ArticleService) getArticlesWithCache(page, size int) (*models.PageResponse, error) {
	cacheKey := s.pageCacheKey(page, size)

	// 1. 查缓存
	cachedData, err := s.redis.Get(cacheKey).Result()
	if err == nil {
		var resp models.PageResponse
		if err := json.Unmarshal([]byte(cachedData), &resp); err != nil {
			return nil, fmt.Errorf("cache unmarshal: %w", err)
		}
		return &resp, nil
	}

	// 2. Redis异常（非key不存在），直接降级查库
	if err != redis.Nil {
		fmt.Printf("warn: redis get %s failed: %v, fallback to db\n", cacheKey, err)
		return s.repo.FindByPage(page, size)
	}

	// 3. 缓存未命中，查库
	resp, err := s.repo.FindByPage(page, size)
	if err != nil {
		return nil, err
	}

	// 4. 回填缓存（失败不影响返回）
	cacheData, _ := json.Marshal(resp)
	if err := s.redis.Set(cacheKey, cacheData, cacheTTL).Err(); err != nil {
		fmt.Printf("warn: redis set %s failed: %v\n", cacheKey, err)
	}

	return resp, nil
}

func (s *ArticleService) pageCacheKey(page, size int) string {
	return fmt.Sprintf("articles:page:%d:size:%d", page, size)
}
