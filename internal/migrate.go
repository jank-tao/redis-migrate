package internal

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

// ======================
// RedisMigrate interface

type IRedisMigrate = *redisMigrate

func NewRedisMigrate(config *MigrateConfig) *redisMigrate {
	return &redisMigrate{
		config: config,
	}
}

// ======================
// RedisMigrate implement

type redisMigrate struct {
	config *MigrateConfig
}

func (m *redisMigrate) MigrateTask(name string) error {
	config := m.config.Tasks[name]
	from, err := m.config.Redis[config.From].Client()
	if err != nil {
		return err
	}
	to, err := m.config.Redis[config.To].Client()
	if err != nil {
		return err
	}

	return Migrate(from, to, config.migrateOption, m.config.Tasks[name].Patterns...)
}

func Migrate(from, to redis.Conn, opt migrateOption, patterns ...string) error {
	dumper := NewRedisDumper(from)
	loader := NewRedisLoader(to)

	for _, pattern := range patterns {
		keys, err := redis.Strings(from.Do("KEYS", pattern))
		if err != nil {
			return err
		}

		for _, key := range keys {
			MigrateKey(dumper, loader, opt, key)
		}
	}

	return nil
}

func MigrateKey(dumper IRedisDumper, loader IRedisLoader, opt migrateOption, key string) error {
	fmt.Println("migrate: ", key)
	obj, err := dumper.Dump(key, opt)
	if err != nil {
		return err
	}
	return loader.Load(obj, opt)
}
