package storage

import (
	"sync"
	"strconv"
)

type ValueType uint8

const (
	RedisString ValueType = iota
	RedisInt
	RedisHash
	RedisList
)

var mutex=&sync.Mutex{}

type RedisClientMock struct{
	keys map[string] int
	intMap map[string] int
	stringMap map[string] interface{}
	listMap map[string] []interface{}
}

func NewMockClient()RedisClientMock{
	keys:=make(map[string] int)
	intMap := make(map[string] int)
	stringMap := make(map[string] interface{})
	listMap :=make(map[string] []interface{})
	return RedisClientMock{
		keys:keys,
		intMap:intMap,
		stringMap:stringMap,
		listMap:listMap,
	}
}

func (rc RedisClientMock) Exists(key string) bool {
	_,ok:=rc.keys[key]
	return ok
}

func (rc RedisClientMock) getList(key string)[]interface{}{
	list,ok:=rc.listMap[key]
	if ok{
		if len(list)== cap(list){
			listnew:=make([]interface{},cap(list),cap(list)+20)
			copy(listnew,list)
			rc.listMap[key]= listnew
			return listnew
		}
		return list
	}else{
		array:=make([]interface{},0,20)
		rc.listMap[key]= array
		return array
	}
}

func (rc RedisClientMock) Lpush(key string,value interface{}) error {
	mutex.Lock()
	list:=rc.getList(key)

	listnew:=make([]interface{},cap(list)+1,cap(list)+20)
	listnew[0]=value
	copy(listnew,list)

	rc.listMap[key]= listnew

	mutex.Unlock()
	return nil
}

func (rc RedisClientMock) Rpush(key string,value interface{}) error {
	mutex.Lock()
	list:=rc.getList(key)
	list=append(list,value)
	rc.listMap[key]=list
	mutex.Unlock()
	return nil
}

func (rc RedisClientMock) Lrange(key string,start,end int)([]interface{}, error) {
	mutex.Lock()
	list:=rc.getList(key)
	result:=list[start:end]
	mutex.Unlock()
	return result,nil
}

func (rc RedisClientMock) Lrem(key string,value interface{},remType int) error {
	return nil
}

func (rc RedisClientMock) Brpop(key string,timeoutSecs int) (interface{},error) {
	return nil,nil
}

func (rc RedisClientMock) Set(key string,value []byte) error {
	rc.stringMap[key]=value
	rc.keys[key]=int(RedisString)
	return nil
}

func (rc RedisClientMock) Get(key string) ([]byte,error) {
	keyType,ok:=rc.keys[key]
	if ok{
		if keyType== int(RedisInt){
			value,ok:=rc.intMap[key]
			if ok{
				return []byte(strconv.Itoa(value)),nil
			}else{
				return nil,nil
			}
		}else{
			value,ok:=rc.stringMap[key]
			if ok{
				return value.([]byte),nil
			}else{
				return nil,nil
			}
		}
	}else{
		return nil,nil
	}
}

func (rc RedisClientMock) Delete(key string) error {
	delete(rc.stringMap,key)
	delete(rc.keys,key)
	return nil
}

func (rc RedisClientMock) Incr(key string,step uint64)(int64, error) {
	mutex.Lock()
	value,ok:=rc.intMap[key]
	var newValue int
	if ok {
		newValue=value+int(step)
	}else{
		newValue=int(step)
	}
	rc.intMap[key]=newValue
	rc.keys[key]=int(RedisInt)
	mutex.Unlock()
	return int64(newValue),nil
}

func (rc RedisClientMock) Decr(key string,step uint64)(int64 ,error ){
	mutex.Lock()
	value,ok:=rc.intMap[key]
	var newValue int
	if ok {
		newValue=value-int(step)
	}else{
		newValue=-int(step)
	}
	rc.intMap[key]=newValue
	mutex.Unlock()
	return int64(newValue),nil
}

func (rc RedisClientMock) MultiGet(keys []interface{})([]interface{},error){
	mutex.Lock()
	var result []interface{}
	for _,key := range keys{
		if value,ok:=rc.stringMap[key.(string)]; ok{
			result=append(result,value)
		}
		rc.keys[key.(string)]=int(RedisString)
	}
	mutex.Unlock()
	return result,nil
}

func (rc RedisClientMock) MultiSet(kvMap map[string][]byte) error {
	mutex.Lock()
	for key,value := range kvMap{
		rc.keys[key]=int(RedisString)
		rc.stringMap[key]=value
	}
	mutex.Unlock()
	return nil
}

func (rc RedisClientMock) ClearAll() error {
	rc.keys=make(map[string] int)
	rc.intMap = make(map[string] int)
	rc.stringMap = make(map[string] interface{})
	rc.listMap =make(map[string] []interface{})
	return nil
}







