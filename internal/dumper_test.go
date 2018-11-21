package internal

import (
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	rc     redis.Conn
	dumper IRedisDumper
)

func init() {
	rc, _ = redis.Dial("tcp", "127.0.0.1:6379")
	dumper = NewRedisDumper(rc)

	rc.Send("MULTI")
	rc.Send("DEL", "dump_test_string", "dump_test_list", "dump_test_hash", "dump_test_set", "dump_test_zset")
	rc.Send("SET", "dump_test_string", "test_value")
	rc.Send("LPUSH", "dump_test_list", "test_value")
	rc.Send("HSET", "dump_test_hash", "test_key", "test_value")
	rc.Send("SADD", "dump_test_set", "test_value")
	rc.Send("ZADD", "dump_test_zset", 42, "test_value")
	rc.Do("EXEC")
}

func Test_DumpString(t *testing.T) {
	ast := assert.New(t)
	obj, err := dumper.Dump("dump_test_string")
	ast.NoError(err)
	ast.NotNil(obj)
	ast.Equal("string", obj.Type)
	ast.Equal([]byte("test_value"), obj.Data)
}

func Test_DumpList(t *testing.T) {
	ast := assert.New(t)
	obj, err := dumper.Dump("dump_test_list")
	ast.NoError(err)
	ast.NotNil(obj)
	ast.Equal("list", obj.Type)
	ast.Equal([]interface{}{[]byte("test_value")}, obj.Data)
}

func Test_DumpHash(t *testing.T) {
	ast := assert.New(t)
	obj, err := dumper.Dump("dump_test_hash")
	ast.NoError(err)
	ast.NotNil(obj)
	ast.Equal("hash", obj.Type)
	data := obj.Data.([]interface{})
	ast.Equal("test_key", string(data[0].([]byte)))
	ast.Equal("test_value", string(data[1].([]byte)))
}

func Test_DumpSet(t *testing.T) {
	ast := assert.New(t)
	obj, err := dumper.Dump("dump_test_set")
	ast.NoError(err)
	ast.NotNil(obj)
	ast.Equal("set", obj.Type)
	ast.Equal([]interface{}{[]byte("test_value")}, obj.Data)
}

func Test_DumpZSet(t *testing.T) {
	ast := assert.New(t)
	obj, err := dumper.Dump("dump_test_zset")
	ast.NoError(err)
	ast.NotNil(obj)
	ast.Equal("zset", obj.Type)
	data := obj.Data.([]interface{})
	ast.Equal("test_value", string(data[0].([]byte)))
	ast.Equal("42", string(data[1].([]byte)))
}
