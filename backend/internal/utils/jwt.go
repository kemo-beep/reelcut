package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func IssueAccessToken(userID, email, secret string, expiry time.Duration) (string, time.Time, error) {
	exp := time.Now().Add(expiry)
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		UserID: userID,
		Email:  email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func IssueRefreshToken(userID, secret string, expiry time.Duration) (string, time.Time, error) {
	exp := time.Now().Add(expiry)
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        uuid.New().String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

func ParseAccessToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func ParseRefreshToken(tokenString, secret string) (userID string, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
}

// ThumbnailClaims is used for short-lived tokens in thumbnail URLs (e.g. in list so <img> can load without Bearer).
type ThumbnailClaims struct {
	jwt.RegisteredClaims
	VideoID string `json:"video_id"`
	UserID  string `json:"user_id"`
}

func IssueThumbnailToken(videoID, userID, secret string, expiry time.Duration) (string, error) {
	exp := time.Now().Add(expiry)
	claims := ThumbnailClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		VideoID: videoID,
		UserID:  userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ParseThumbnailToken(tokenString, secret string) (*ThumbnailClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ThumbnailClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*ThumbnailClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// OneTimeClaims for password reset and email verification tokens.
type OneTimeClaims struct {
	jwt.RegisteredClaims
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Purpose string `json:"purpose"` // "password_reset" or "email_verify"
}

func IssuePasswordResetToken(userID, email, secret string, expiry time.Duration) (string, error) {
	exp := time.Now().Add(expiry)
	claims := OneTimeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		UserID:  userID,
		Email:   email,
		Purpose: "password_reset",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func IssueEmailVerifyToken(userID, email, secret string, expiry time.Duration) (string, error) {
	exp := time.Now().Add(expiry)
	claims := OneTimeClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		UserID:  userID,
		Email:   email,
		Purpose: "email_verify",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseOneTimeToken(tokenString, secret, purpose string) (*OneTimeClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &OneTimeClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*OneTimeClaims)
	if !ok || !token.Valid || claims.Purpose != purpose {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
