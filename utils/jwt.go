package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"lanshan_chat/app/api/global"
	"time"
)

type MyClaims struct {
	Uid int64
	jwt.RegisteredClaims
}

type CustomClaims struct {
	ID interface{}
	jwt.RegisteredClaims
}

func GenCustomToken(id interface{}, expiresTime int64) (string, error) {
	claim := CustomClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Truncate(time.Second)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresTime) * time.Second)),
		},
	}
	return genCustomTokenWithClaim(claim)
}

func genCustomTokenWithClaim(claim jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(global.Config.AuthConfig.JwtConfig.SecretKey))
}

func GenToken(uid int64) (string, error) {
	claim := MyClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			NotBefore: jwt.NewNumericDate(time.Now().Truncate(time.Second)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.Config.AuthConfig.JwtConfig.ExpiresTime) * time.Second)),
		},
	}
	return genCustomTokenWithClaim(claim)
}

func ParseToken(tokenStr string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&MyClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(global.Config.AuthConfig.JwtConfig.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*MyClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid Token")
}

func ParseCustomToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(global.Config.AuthConfig.JwtConfig.SecretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid Token")
}
