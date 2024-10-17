package jwt

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	UserID   int64  `json:"userid"`
	Username string `json:"username"`
}

var jwtKey []byte

// Initialize jwtKey from Viper configuration
func InitJWT() error {
	jwtKey = []byte(viper.GetString("JWT_SECRET"))
	if len(jwtKey) == 0 {
		return errors.New("JWT secret not found in configuration")
	}
	return nil
}

func GenerateToken(userID int64, userName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid":   userID,
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return 0, errors.New("token expired")
		}
		if userID, ok := claims["userid"].(float64); ok {
			return int64(userID), nil
		}
		return 0, errors.New("userid claim is missing or not a float64")
	}

	return 0, errors.New("invalid token")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func DecodeToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := int64(claims["userid"].(float64))
		username := claims["username"].(string)
		return &TokenClaims{
			UserID:   userID,
			Username: username,
		}, nil
	}

	return nil, errors.New("invalid token")
}

func GetUserInfo(c *gin.Context) (*TokenClaims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is missing")
	}

	tokenString := authHeader[len("Bearer "):]

	tokenClaims, err := DecodeToken(tokenString)
	if err != nil {
		return nil, err
	}

	return tokenClaims, nil
}
