package storage

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

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

	_, err := conn.Do("LPUSH", key, value)
	return err
}

func (rc RedisClient) Rpush(key string,value interface{}) error {
	conn:=rc.connectInit()
	defer conn.Close()

	_, err := conn.Do("RPUSH", key, value)
	return err
}

func (rc RedisClient) Lrange(key string,start,end int)([]interface{}, error) {
	conn:=rc.connectInit()
	defer conn.Close()

	v, err := conn.Do("LRANGE", key, start,end)
	return v.([]interface{}),err
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
