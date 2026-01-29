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

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1ZCI6eyJVc2VySWQiOjF9LCJidWZmVGltZSI6MTc3MDE1MzAxMywiZXhwIjoxNzY5NzM1NDEzLCJpYXQiOjE3Njk2OTIyMTN9.4JREwP-vkrcG4Rqb2s1dQiEpd6eaJGMHWQ2fTLMiN7k
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVZCI6eyJVc2VySWQiOjF9LCJCdWZmVGltZSI6MTc3MDA4MjkyMiwiZXhwIjoxNzY5NjY1MzIyLCJpYXQiOjE3Njk2MjIxMjJ9.mf3iea4YCHTH1u3ycIt_cNJ2DXTZMOtUPGbr3ogLD5Q
