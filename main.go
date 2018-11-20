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
	configFilepath := flag.String("c", "example.config.json", "config file")
	taskName := flag.String("t", "default", "task name")
	flag.Parse()

	config := new(internal.MigrateConfig)
	bytes, err := ioutil.ReadFile(*configFilepath)
	check(err)
	json.Unmarshal(bytes, config)

	m := internal.NewRedisMigrate(config)
	check(m.RunTask(*taskName))
}
