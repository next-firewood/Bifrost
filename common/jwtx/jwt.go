package jwtx

import (
	"bifrost/common/errorx"
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
	Ud       UserData         `json:"ud"`
	BuffTime *jwt.NumericDate `json:"buff_time,omitempty"`
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
func (a Auth) ValidateToken(tokenStr string) (*CustomClaims, error) {
	cla := &CustomClaims{}

	parser := jwt.NewParser(
		jwt.WithoutClaimsValidation(),
	)

	token, err := parser.ParseWithClaims(
		tokenStr,
		cla,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.AccessSecret), nil
		},
	)
	if err != nil {
		return nil, errorx.NewCodeError(errorx.TokenExpired)
	}

	if !token.Valid {
		return nil, errorx.NewCodeError(errorx.TokenExpired)
	}

	now := time.Now()

	if cla.ExpiresAt == nil {
		return nil, errorx.NewCodeError(errorx.TokenExpired)
	}

	if now.Before(cla.ExpiresAt.Time) {
		return cla, nil
	}

	if cla.BuffTime != nil && now.Before(cla.BuffTime.Time) {
		return cla, errorx.NewCodeError(errorx.TokenRefresh)
	}

	return nil, errorx.NewCodeError(errorx.TokenExpired)
}
