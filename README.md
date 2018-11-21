# redis-migrate
redis-migrate is a tool to migrate data between redis databases.

## install
```
go get -u -v https://github.com/liguangsheng/redis-migrate
```

usage:
```
redis-migrate -c `/path/to/config/file` -t `task_name`
```

example:
```
redis-migrate -c redis-migrate.json -t default
```

## example config
```
{
  "redis": {
    "remote_redis": {
      "addr": "10.0.0.1:6379",
      "password": "password",
      "db": 1,
      "ssh_config": {
        "addr": "123.124.125.126:22",
        "private_key_path": "/path/to/your/.ssh/id_rsa",
        "username": "username"
      }
    },
    "local_redis": {
      "addr": "127.0.0.1:6379",
      "password": "",
      "db": 1,
      "ssh_config": null
    }
  },
  "tasks": {
    "default": {
      "from": "remote_redis",
      "to": "local_redis",
      "patterns": [
        "*"
      ]
    }
  }
}
```


```
$ redis-migrate --help
Usage of redis-migrate:
  -c string
    	config file (default "example.config.json")
  -t string
    	task name (default "default")
```

# TODO

- dump to file
- load from file