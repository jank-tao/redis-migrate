# redis-migrate
A tool to migrate data bwtween redis servers

redis 数据迁移工具

```
go get -v https://github.com/liguangsheng/redis-migrate
```
```
redis-migrate -c redis-migrate.json -t task1
```

# example config
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

- pipeline加速
- goroutine加速
- 更多数据类型支持
- 文件导入导出