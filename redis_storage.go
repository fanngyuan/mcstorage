package storage

import (
	"reflect"
	"encoding/json"
)

type RedisStorage struct {
	client            RedisClient
	KeyPrefix         string
	DefaultExpireTime int
	T                 reflect.Type
}

func (this RedisStorage) Get(key interface{}) (interface{}, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	data, err := this.client.Get(cacheKey)
	object, err := bytesToInterface(data,this.T)
	if err != nil {
		return nil, err
	}
	return object,nil

}

func (this RedisStorage) Set(key interface{}, object interface{}) error {
	buf, err := json.Marshal(object)
	if err != nil {
		return err
	}
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Set(keyCache,buf)
	return nil
}

func (this RedisStorage) MultiGet(keys []interface{}) (map[interface{}]interface{}, error){
	return nil,nil
}

func (this RedisStorage) MultiSet(map[interface{}]interface{}) error{
	return nil
}

func (this RedisStorage) Delete(key interface{}) error{
	return nil
}
