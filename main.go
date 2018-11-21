package main

import (
	"encoding/json"
	"flag"
	"github.com/liguangsheng/redis-migrate/internal"
	"io/ioutil"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// parse cli args
	configFilepath := flag.String("c", "redis-migrate.json", "path to config file")
	taskName := flag.String("t", "default", "task name")
	flag.Parse()

	// load config
	config := new(internal.MigrateConfig)
	bytes, err := ioutil.ReadFile(*configFilepath)
	check(err)
	check(json.Unmarshal(bytes, config))

	// do migrate
	check(internal.NewRedisMigrate(config).MigrateTask(*taskName))
}
