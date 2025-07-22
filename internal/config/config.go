package config

import "time"

var (
	JWTSecret      = []byte("your-secret-key")
	JWTTokenExpiry = time.Hour * 24
)
