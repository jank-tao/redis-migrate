package internal

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// =====================
// RedisDumper interface
type IRedisDumper interface {
	Dump(key string, opt migrateOption) (*RedisObject, error)
}

func NewRedisDumper(rc redis.Conn) IRedisDumper {
	return &redisDumper{
		rc: rc,
	}
}

// =====================
// RedisDumper implement
type redisDumper struct {
	rc redis.Conn
}

func (d *redisDumper) Dump(key string, opt migrateOption) (*RedisObject, error) {
	d.rc.Send("EXISTS", key)
	d.rc.Send("TYPE", key)
	d.rc.Send("TTL", key)
	d.rc.Flush()

	exists, err := redis.Bool(d.rc.Receive())
	if err != nil || !exists {
		return nil, err
	}

	_type, err := redis.String(d.rc.Receive())
	if err != nil {
		return nil, err
	}

	ttl, err := redis.Int(d.rc.Receive())
	if err != nil || (-1 < ttl && ttl < 5) {
		return nil, err
	}

	var data interface{}
	switch _type {
	case "string":
		data, err = d.dumpString(key)
	case "list":
		data, err = d.dumpList(key)
	case "hash":
		data, err = d.dumpHash(key)
	case "set":
		data, err = d.dumpSet(key)
	case "zset":
		data, err = d.dumpZSet(key)
	default:
		return nil, errors.New(fmt.Sprintf("not support type: %s", _type))
	}

	if err != nil {
		return nil, err
	}

	return &RedisObject{
		Key:  key,
		Type: _type,
		TTL:  ttl,
		Data: data,
	}, nil
}

func (d *redisDumper) dumpString(key string) (interface{}, error) {
	return d.rc.Do("GET", key)
}

func (d *redisDumper) dumpList(key string) (interface{}, error) {
	return d.rc.Do("LRANGE", key, 0, -1)
}

func (d *redisDumper) dumpHash(key string) (interface{}, error) {
	return d.rc.Do("HGETALL", key)
}

func (d *redisDumper) dumpSet(key string) (interface{}, error) {
	return d.rc.Do("SMEMBERS", key)
}

func (d *redisDumper) dumpZSet(key string) (interface{}, error) {
	return d.rc.Do("ZRANGE", key, 0, -1, "WITHSCORES")
}
