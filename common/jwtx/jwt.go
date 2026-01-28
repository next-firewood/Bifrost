package jwtx

import (
	"bifrost/common/errorx"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Auth struct {
	AccessSecret string `json:"accessSecret" yaml:"accessSecret"`
	AccessExpire string `json:"accessExpire" yaml:"accessExpire"`
	BuffTime     string `json:"buffTime" yaml:"buffTime"`
}

type UserData map[string]interface{}

type CustomClaims struct {
	Ud UserData

	BuffTime *jwt.NumericDate
	jwt.RegisteredClaims
}

func (a Auth) GenToken(c *CustomClaims) (string, error) {
	now := time.Now()

	expDur, err := time.ParseDuration(a.AccessExpire)
	if err != nil {
		return "", err
	}

	bufDur, err := time.ParseDuration(a.BuffTime)
	if err != nil {
		return "", err
	}

	c.ExpiresAt = jwt.NewNumericDate(now.Add(expDur))
	c.IssuedAt = jwt.NewNumericDate(now)
	c.BuffTime = jwt.NewNumericDate(now.Add(bufDur))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return token.SignedString([]byte(a.AccessSecret))
}
func (a Auth) ValidateToken(tokenStr string) (cla *CustomClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.AccessSecret), nil
	})
	if err != nil && !errors.Is(err, jwt.ErrTokenExpired) {
		return nil, errorx.NewCodeError(errorx.TokenExpired)
	}

	if errors.Is(err, jwt.ErrSignatureInvalid) {
		return nil, errorx.NewCodeError(errorx.TokenExpired)
	}

	c, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errorx.NewCodeError(errorx.TokenExpired)
	}

	if errors.Is(err, jwt.ErrTokenExpired) {
		if c.BuffTime.After(time.Now()) {
			return c, errorx.NewCodeError(errorx.TokenRefresh)
		} else {
			return nil, errorx.NewCodeError(errorx.TokenExpired)
		}
	}

	if !token.Valid {
		return nil, errorx.NewCodeError(errorx.TokenExpired)
	}

	return c, nil
}
