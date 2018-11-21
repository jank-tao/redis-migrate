package internal

import (
	"github.com/gomodule/redigo/redis"
)

// =====================
// RedisLoader interface

type IRedisLoader interface {
	Load(obj *RedisObject, opt migrateOption) error
}

func NewRedisLoader(rc redis.Conn) IRedisLoader {
	return &redisLoader{
		rc: rc,
	}
}

// =====================
// RedisLoader implement
type redisLoader struct {
	rc redis.Conn
}

func (l *redisLoader) Load(obj *RedisObject, opt migrateOption) error {
	if obj == nil {
		return nil
	}

	exists, err := redis.Bool(l.rc.Do("EXISTS", obj.Key))
	if err != nil {
		return err
	}
	if exists {
		l.rc.Do("DEL", obj.Key)
	}

	switch obj.Type {
	case "string":
		err = l.loadString(obj)
	case "list":
		err = l.loadList(obj)
	case "hash":
		err = l.loadHash(obj)
	case "set":
		err = l.loadSet(obj)
	case "zset":
		err = l.loadZSet(obj)
	}

	if err != nil {
		return err
	}

	if !opt.IgnoreTTL && obj.TTL != -1 {
		l.rc.Do("EXPIRE", obj.Key, obj.TTL)
	}

	return nil
}

func (l *redisLoader) loadString(obj *RedisObject) error {
	_, err := l.rc.Do("SET", obj.Key, obj.Data)
	return err
}

func (l *redisLoader) loadList(obj *RedisObject) error {
	values, err := redis.Values(obj.Data, nil)
	if err != nil {
		return err
	}
	args := append([]interface{}{obj.Key}, values...)
	_, err = l.rc.Do("LPUSH", args...)
	return err
}

func (l *redisLoader) loadHash(obj *RedisObject) error {
	values, err := redis.Values(obj.Data, nil)
	if err != nil {
		return err
	}
	args := append([]interface{}{obj.Key}, values...)
	_, err = l.rc.Do("HMSET", args...)
	return err
}

func (l *redisLoader) loadSet(obj *RedisObject) error {
	values, err := redis.Values(obj.Data, nil)
	if err != nil {
		return err
	}
	args := append([]interface{}{obj.Key}, values...)
	_, err = l.rc.Do("SADD", args...)
	return err
}

func (l *redisLoader) loadZSet(obj *RedisObject) error {
	values, err := redis.Values(obj.Data, nil)
	if err != nil {
		return err
	}

	args := []interface{}{obj.Key}
	for i := 0; i < len(values); i += 2 {
		member, _ := redis.String(values[i], nil)
		score, _ := redis.Float64(values[i+1], nil)
		args = append(args, score, member)
	}

	_, err = l.rc.Do("ZADD", args...)
	return err
}
