package oauth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type StateClaims struct {
	Provider string `json:"provider"`
	jwt.RegisteredClaims
}

func GenerateState(secret string, provider string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := StateClaims{
		Provider: provider,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseState(rawState string, secret string, provider string) (*StateClaims, error) {
	parsed, err := jwt.ParseWithClaims(rawState, &StateClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*StateClaims)
	if !ok || !parsed.Valid || claims.Provider != provider {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
