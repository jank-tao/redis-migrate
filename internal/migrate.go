package internal

import (
	"fmt"
	"gopkg.in/redis.v5"
)

type RedisMigrate struct {
	config *MigrateConfig
}

func NewRedisMigrate(config *MigrateConfig) *RedisMigrate {
	return &RedisMigrate{
		config: config,
	}
}

func (m *RedisMigrate) RunTask(name string) error {
	fromclient, err := m.config.Redis[m.config.Tasks[name].From].Client()
	if err != nil {
		return err
	}
	toclient, err := m.config.Redis[m.config.Tasks[name].To].Client()
	if err != nil {
		return err
	}

	runner := &MigrateRunner{
		fromClient: fromclient,
		toClient:   toclient,
		dumper:     &RedisDumper{fromclient},
		loader:     &RedisLoader{toclient},
	}
	return runner.Migrate(m.config.Tasks[name].Patterns...)
}

type MigrateRunner struct {
	fromClient *redis.Client
	toClient   *redis.Client
	dumper     *RedisDumper
	loader     *RedisLoader
}

func (m *MigrateRunner) Migrate(patterns ...string) error {
	for _, pattern := range patterns {
		keys, err := m.fromClient.Keys(pattern).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			if err := m.MigrateKey(key); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *MigrateRunner) MigrateKey(key string) error {
	fmt.Println("migrate: ", key)
	item, err := m.dumper.Dump(key)
	if err != nil {
		return err
	}
	return m.loader.Load(item)
}
