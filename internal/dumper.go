package internal

import (
	"gopkg.in/redis.v5"
)

type RedisDumper struct {
	Client *redis.Client
}

func (d *RedisDumper) DumpKeys(patterns ...string) ([]*RedisItem, error) {
	res := make([]*RedisItem, 0)
	for _, pattern := range patterns {
		keys, err := d.Client.Keys(pattern).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			item, err := d.Dump(key)
			if err != nil {
				return nil, err
			}
			res = append(res, item)
		}
	}
	return res, nil
}

func (d *RedisDumper) Dump(key string) (*RedisItem, error) {
	keyType, err := d.Client.Type(key).Result()
	if err != nil {
		return nil, err
	}

	switch keyType {
	case "string":
		return d.dumpString(key)
	case "list":
		return d.dumpList(key)
	case "hash":
		return d.dumpHash(key)
	case "set":
		return d.dumpSet(key)
	}
	return nil, nil
}

func (d *RedisDumper) dumpString(key string) (*RedisItem, error) {
	data, err := d.Client.Get(key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	ttl, err := d.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &RedisItem{
		Key:        key,
		Type:       "string",
		TTL:        ttl,
		StringData: data,
	}, nil
}

func (d *RedisDumper) dumpList(key string) (*RedisItem, error) {
	data, err := d.Client.LRange(key, 0, -1).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	ttl, err := d.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &RedisItem{
		Key:      key,
		Type:     "list",
		TTL:      ttl,
		ListData: data,
	}, nil
}

func (d *RedisDumper) dumpHash(key string) (*RedisItem, error) {
	data, err := d.Client.HGetAll(key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	ttl, err := d.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &RedisItem{
		Key:      key,
		Type:     "hash",
		TTL:      ttl,
		HashData: data,
	}, nil
}

func (d *RedisDumper) dumpSet(key string) (*RedisItem, error) {
	data, err := d.Client.SMembers(key).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	ttl, err := d.Client.TTL(key).Result()
	if err != nil {
		return nil, err
	}

	return &RedisItem{
		Key:      key,
		Type:     "set",
		TTL:      ttl,
		ListData: data,
	}, nil
}