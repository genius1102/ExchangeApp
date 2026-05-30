package services

import (
	"errors"

	"exchangeapp/backend/models"
	"exchangeapp/backend/repository"
	"exchangeapp/backend/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register 注册：密码加密 → JWT生成 → 入库，返回token
func (s *AuthService) Register(user *models.User) (string, error) {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return "", err
	}
	user.Password = hash

	token, err := utils.GenerateJWT(user.UserName)
	if err != nil {
		return "", err
	}

	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	return token, nil
}

// Login 登录：查找用户 → 验证密码 → 生成JWT，返回token
func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		// 将ORM错误转换为业务错误，让Controller不感知数据库层
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utils.ErrUserNotFound
		}
		return "", err
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", utils.ErrInvalidPassword
	}

	token, err := utils.GenerateJWT(user.UserName)
	if err != nil {
		return "", err
	}

	return token, nil
}
