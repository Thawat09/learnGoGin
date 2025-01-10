package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId         int    `json:"userId"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	RoleId         int    `json:"roleId"`
	RoleName       string `json:"roleName"`
	DepartmentId   int    `json:"departmentId"`
	DepartmentName string `json:"departmentName"`
	jwt.RegisteredClaims
}

func Encrypt(message string) (string, error) {
	secretKey := os.Getenv("SECRETKEYDATA")

	if secretKey == "" {
		return "", errors.New("missing SECRETKEYDATA in environment variables")
	}

	hash := sha256.Sum256([]byte(secretKey))
	block, err := aes.NewCipher(hash[:])

	if err != nil {
		return "", err
	}

	iv := []byte("123456789012")
	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, iv, []byte(message), nil)
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func Decrypt(encodedMessage string) (string, error) {
	secretKey := os.Getenv("SECRETKEYDATA")

	if secretKey == "" {
		return "", errors.New("missing SECRETKEYDATA in environment variables")
	}

	hash := sha256.Sum256([]byte(secretKey))
	block, err := aes.NewCipher(hash[:])

	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encodedMessage)

	if err != nil {
		return "", err
	}

	iv := []byte("123456789012")

	if len(ciphertext) < len(iv) {
		return "", errors.New("ciphertext too short")
	}

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)

	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func ParseToken(tokenString string) (*Claims, error) {
	secretKey := os.Getenv("SECRETTOKENKEY")

	if secretKey == "" {
		return nil, errors.New("missing SECRETTOKENKEY in environment variables")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}

func ParseRefeshToken(tokenString string) (*Claims, error) {
	secretKey := os.Getenv("SECRETREFRESHTOKENKEY")

	if secretKey == "" {
		return nil, errors.New("missing SECRETREFRESHTOKENKEY in environment variables")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}
