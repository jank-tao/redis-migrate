package internal

type RedisObject struct {
	Key  string      `json:"key"`
	Type string      `json:"type"`
	TTL  int         `json:"ttl"`
	Data interface{} `json:"data"`
}
