package storage

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"reflect"
)

type Redis interface{

	Exists(key string) bool

	Lpush(key string,value interface{}) error

	Rpush(key string,value interface{}) error

	Lrange(key string,start,end int)([]interface{}, error)

	Lrem(key string,value interface{},remType int) error

	Brpop(key string,timeoutSecs int) (interface{},error)

	Set(key string,value []byte) error

	Get(key string) ([]byte,error)

	Delete(key string) error

	Incr(key string,step uint64)(int64, error)

	Decr(key string,step uint64)(int64 ,error )

	MultiGet(keys []interface{})([]interface{},error)

	MultiSet(kvMap map[string][]byte) error

	ClearAll() error

}

type RedisClient struct {
	pool     *redis.Pool
	addr     string
}

func (rc RedisClient) Exists(key string) bool {
	conn:=rc.connectInit()
	defer conn.Close()
	v, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return v
}

func (rc RedisClient) Lpush(key string,value interface{}) error {
	conn:=rc.connectInit()
	defer conn.Close()

	if reflect.TypeOf(value).Kind()==reflect.Slice{
		s := reflect.ValueOf(value)
		values:=make([]interface{},s.Len()+1)
		values[0]=key
        for i := 1; i <= s.Len(); i++ {
			values[i]=s.Index(i-1).Interface()
        }
		_, err := conn.Do("LPUSH", values...)
		return err
	}else{
		_, err := conn.Do("LPUSH", key, value)
		return err
	}
}

func (rc RedisClient) Rpush(key string,value interface{}) error {
	conn:=rc.connectInit()
	defer conn.Close()

	if reflect.TypeOf(value).Kind()==reflect.Slice{
		s := reflect.ValueOf(value)
		values:=make([]interface{},s.Len()+1)
		values[0]=key
        for i := 1; i <= s.Len(); i++ {
			values[i]=s.Index(i-1).Interface()
        }
		_, err := conn.Do("RPUSH", values...)
		return err
	}else{
		_, err := conn.Do("RPUSH", key, value)
		return err
	}
}

func (rc RedisClient) Lrange(key string,start,end int)([]interface{}, error) {
	conn:=rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("LRANGE", key, start,end)
	return v.([]interface{}),err
}

func (rc RedisClient) Lrem(key string,value interface{},remType int) error {
	conn:=rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("LREM", key,remType ,value)
	return err
}

func (rc RedisClient) Brpop(key string,timeoutSecs int) (interface{},error) {
	conn:=rc.connectInit()
	defer conn.Close()

	var val interface{}
	var err error
	if timeoutSecs<0{
		val, err = conn.Do("BRPOP", key,0)
	}else{
		val, err = conn.Do("BRPOP", key,timeoutSecs)
	}
	values,err:=redis.Values(val,err)
	if err!=nil{
		return nil,err
	}
	return string(values[1].([]byte)),err
}

func (rc RedisClient) Set(key string,value []byte) error {
	conn:=rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("SET", key,value)
	return err
}

func (rc RedisClient) Get(key string) ([]byte,error) {
	conn:=rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("GET", key)
	if (err!=nil||v==nil){
		return nil,err
	}
	return v.([]byte),err
}

func (rc RedisClient) Delete(key string) error {
	conn:=rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func (rc RedisClient) Incr(key string,step uint64)(int64, error) {
	conn:=rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("INCRBY", key,step)
	if err!=nil{
		return 0,nil
	}
	return value.(int64),err
}

func (rc RedisClient) Decr(key string,step uint64)(int64 ,error ){
	conn:=rc.connectInit()
	defer conn.Close()

	value, err := conn.Do("DECRBY", key,step)
	if err!=nil{
		return 0,nil
	}
	return value.(int64),err
}

func (rc RedisClient) MultiGet(keys []interface{})([]interface{},error){
	conn:=rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("MGET", keys...)
	return v.([]interface{}),err
}

func (rc RedisClient) MultiSet(kvMap map[string][]byte) error {
	conn:=rc.connectInit()
	defer conn.Close()

	var values []interface{}
	for key,value :=range(kvMap){
		values=append(values,key)
		values=append(values,value)
	}

	_, err := conn.Do("MSET", values...)
	return err
}

func (rc RedisClient) ClearAll() error {
	conn:=rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("FLUSHALL")
	return err
}

func (rc RedisClient) connectInit() redis.Conn {
	conn:=rc.pool.Get()
	return conn
}

func InitClient(addr string) (RedisClient,error) {
	pool := &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return RedisClient{pool,addr},nil
}
