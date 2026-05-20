package utils

import (
	"time"
	"errors"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecret string

// 密码加密
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// JWT生成
func GenerateJWT(username string) (string,error){
	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":  time.Now().Add(72 * time.Hour).Unix(),
	})

	// 签名JWT
	signedToken, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return "Bearer " + signedToken, nil

}

// 密码验证
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// JWT解析
func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	// 解析签名
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, errors.New("unexpected Signing Method")
		}
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return "", err
	}

	// 验证JWT并提取用户名
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		username, ok := claims["username"].(string)
		if !ok{
			return "", errors.New("username claim is not a string")
		}
		return username, nil
	}

	return "", err
}



