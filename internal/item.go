package internal

import (
	"gopkg.in/redis.v5"
	"time"
)

type RedisItem struct {
	Key        string
	Type       string
	TTL        time.Duration
	StringData string
	ListData   []string
	HashData   map[string]string
	ZSetData   []redis.Z
}
