package internal

import "time"

type RedisItem struct {
	Key        string
	Type       string
	TTL        time.Duration
	StringData string
	ListData   []string
	HashData   map[string]string
}
