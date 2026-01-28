package jwtx

import (
	"testing"
)

func TestJwt(t *testing.T) {
	auth := Auth{
		AccessExpire: "12h",
		AccessSecret: "reeee",
		BuffTime:     "128h",
	}

	token, err := auth.GenToken(&CustomClaims{
		Ud: map[string]interface{}{
			"UserId": 1,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(token)
}
