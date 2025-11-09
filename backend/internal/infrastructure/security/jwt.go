package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPair struct {
	AccessToken  string `json:"AccessToken"`
	RefreshToken string `json:"RefreshToken"`
}

type JWTService struct {
	secret []byte
	ttl    time.Duration 
}

func NewJWTService(secret string, ttl time.Duration) *JWTService {
	return &JWTService{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func NewJWTServiceFromEnv() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	return NewJWTService(secret, 15*time.Minute)
}

func (j *JWTService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
}


func (j *JWTService) GenerateTokenPair(userID uint, role string) (TokenPair, string, error) {
	tokenID := generateTokenID()

	accessClaims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"tid":  tokenID,
		"exp":  time.Now().Add(j.ttl).Unix(),
		"iat":  time.Now().Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"sub": userID,
		"tid": tokenID,
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 días
	}

	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessStr, err := access.SignedString(j.secret)
	if err != nil {
		return TokenPair{}, "", err
	}

	refreshStr, err := refresh.SignedString(j.secret)
	if err != nil {
		return TokenPair{}, "", err
	}

	return TokenPair{
		AccessToken:  accessStr,
		RefreshToken: refreshStr,
	}, tokenID, nil
}

func generateTokenID() string {
	// para dev.  más robusto github.com/google/uuid
	return time.Now().Format("20060102150405.000000000")
}
