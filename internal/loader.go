package internal

import (
	"gopkg.in/redis.v5"
)

type RedisLoader struct {
	Client *redis.Client
}

func (l *RedisLoader) LoadItems(items ...*RedisItem) error {
	for _, item := range items {
		if err := l.Load(item); err != nil {
			return err
		}
	}
	return nil
}

func (l *RedisLoader) Load(item *RedisItem) error {
	if item == nil {
		return nil
	}

	ok, err := l.Client.Exists(item.Key).Result()
	if err != nil {
		return err
	}

	if ok {
		l.Client.Del(item.Key)
	}

	switch item.Type {
	case "string":
		return l.loadString(item)
	case "list":
		return l.loadList(item)
	case "hash":
		return l.loadHash(item)
	case "set":
		return l.loadSet(item)
	case "zset":
		return l.loadZSet(item)
	}
	return nil
}

func (l *RedisLoader) loadString(item *RedisItem) error {
	if err := l.Client.Set(item.Key, item.StringData, -1).Err(); err != nil {
		return err
	}

	//if item.TTL != time.Second*-1 {
	//	l.Client.Expire(item.Key, item.TTL)
	//}
	return nil
}

func (l *RedisLoader) loadList(item *RedisItem) error {
	pipe := l.Client.Pipeline()
	for _, i := range item.ListData {
		pipe.LPush(item.Key, i)
	}

	//if item.TTL != time.Second*-1 {
	//	pipe.Expire(item.Key, item.TTL)
	//}
	pipe.Exec()
	return nil
}

func (l *RedisLoader) loadHash(item *RedisItem) error {
	if err := l.Client.HMSet(item.Key, item.HashData).Err(); err != nil {
		return err
	}

	//if item.TTL != time.Second*-1 {
	//	l.Client.Expire(item.Key, item.TTL)
	//}
	return nil
}

func (l *RedisLoader) loadSet(item *RedisItem) error {
	pipe := l.Client.Pipeline()
	for _, i := range item.ListData {
		pipe.SAdd(item.Key, i)
	}

	//if item.TTL != time.Second*-1 {
	//	pipe.Expire(item.Key, item.TTL)
	//}
	pipe.Exec()
	return nil
}

func (l *RedisLoader) loadZSet(item *RedisItem) error {
	pipe := l.Client.Pipeline()
	for _, i := range item.ZSetData {
		pipe.ZAdd(item.Key, i)
	}

	//if item.TTL != time.Second*-1 {
	//	pipe.Expire(item.Key, item.TTL)
	//}
	pipe.Exec()
	return nil
}
