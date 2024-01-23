package config

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var Redis = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("redis_host"),
	Password: "",
	// Password: "Qonics!",
	// Password: os.Getenv("redis_password"),
	DB: 0,
})
