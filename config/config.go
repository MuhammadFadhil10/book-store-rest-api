package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("2312312j3hkj213n23k1jn1ku21nh1k2nlk21jn3lk21n3")

type JWTClaim struct {
	UserName string
	jwt.RegisteredClaims
}