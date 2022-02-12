package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret            = getenv("JWT_SECRET", "randomchatjwtcred")
	jwtExpirationSeconds int64
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

func init() {
	var err error
	jwtExpirationSeconds, err = strconv.ParseInt(getenv("JWT_EXPIRATION_SECONDS", "86400"), 10, 0)
	if err != nil {
		panic(err)
	}
}

type JWTClaims struct {
	ChannelID uint64
	jwt.StandardClaims
}

type AuthPayload struct {
	AccessToken string
}

type AuthResponse struct {
	ChannelID uint64
	Expired   bool
}

func Auth(authPayload *AuthPayload) (*AuthResponse, error) {
	token, err := parseToken(authPayload.AccessToken)
	if err != nil {
		v := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return &AuthResponse{
				Expired: true,
			}, nil
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !(ok && token.Valid) {
		return nil, ErrInvalidToken
	}

	return &AuthResponse{
		ChannelID: claims.ChannelID,
		Expired:   false,
	}, nil
}

func newJWT(channelID uint64) (string, error) {
	expiredAt := time.Now().Add(time.Duration(jwtExpirationSeconds) * time.Second).Unix()
	jwtClaims := &JWTClaims{
		ChannelID: channelID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

func parseToken(accessToken string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(accessToken, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})
}
