package storage

import (
	//"github.com/bradfitz/gomemcache/memcache"
)

type RedisListStorage struct {
	RedisStorage
}

func (this RedisListStorage) Get(key interface{}) (interface{}, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	data, err := this.client.Lrange(cacheKey,0,MAXLEN)
	return data,err
}

func (this RedisListStorage) Set(key interface{}, object interface{}) error {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	return this.client.Rpush(cacheKey,object)
}

func (this RedisListStorage) MultiGet(keys []interface{}) (map[interface{}]interface{}, error){
	result:=make(map[interface{}] interface{})
	for key :=range keys{
		value,err:=this.Get(key)
		if err!=nil{
			continue
		}
		result[key]=value
	}
	return result,nil
}

func (this RedisListStorage) MultiSet(valueMap map[interface{}]interface{}) error{
	for key,value :=range valueMap{
		cacheKey,err:=BuildCacheKey(this.KeyPrefix, key)
		if err!=nil{
			continue
		}
		this.client.Lpush(cacheKey,value)
	}
	return nil
}

func (this RedisListStorage) Getlimit(key,sinceId,maxId interface{},page,count int)(interface{},error){
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	if sinceId==0 && maxId==0{
		return this.client.Lrange(cacheKey,(page-1)*count,page*count-1)
	}
	obj,err:=this.Get(key)
	if err!=nil{
		return nil,err
	}
	return Page(obj.(Pagerable),sinceId,maxId,page,count),nil
}

func (this RedisListStorage) AddItem(key interface{},item interface{})error{
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	return this.client.Lpush(cacheKey,item)
}

func (this RedisListStorage) DeleteItem(key interface{},item interface{})error{
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	return this.client.Lrem(cacheKey,item,0)
}
