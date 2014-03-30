package storage

import (
	"strconv"
	"sort"
	"reflect"
)

type RedisListStorage struct {
	RedisStorage
	decodeList DecodeList
}

type DecodeList func(data []interface{})Pagerable

func NewRedisListStorage(serverUrl string, keyPrefix string, defaultExpireTime int,decode DecodeList)(RedisListStorage,error){
	client,err:=InitClient(serverUrl)
	redisStorage:=RedisStorage{client, keyPrefix, defaultExpireTime, nil}
	return RedisListStorage{redisStorage,decode},err
}

func (this RedisListStorage) Get(key interface{}) (interface{}, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	data, err := this.client.Lrange(cacheKey,0,MAXLEN)
	if err != nil {
		return nil, err
	}
	pageData:=this.decodeList(data)
	return pageData,err
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

	if reflect.ValueOf(sinceId).Int()==0 && reflect.ValueOf(maxId).Int()==0{
		data,err:=this.client.Lrange(cacheKey,(page-1)*count,page*count-1)
		if err!=nil{
			return nil,err
		}
		pageData:=this.decodeList(data)
		return pageData,nil
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

func DecodeIntReversedSlice(data []interface{})Pagerable{
	intArray:= make([]int,len(data))
	for i,item := range data{
		intItem,err:=strconv.Atoi(string(item.([]byte)))
		if err!=nil{
			continue
		}
		intArray[i]=intItem
	}
	sort.Sort(sort.Reverse(sort.IntSlice(intArray)))
	slice:=IntReversedSlice(intArray)
	return slice
}
