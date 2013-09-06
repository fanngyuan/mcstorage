package storage

import (
	"bytes"
	"encoding/json"
	"github.com/golang/groupcache"
	"reflect"
)

type GroupCacheKvStorage struct{
	CacheGroup *groupcache.Group
	DefaultExpireTime int
	T          reflect.Type
}

func (this *GroupCacheKvStorage) Get(key interface{}) (interface{}, error) {
	var data []byte
	this.CacheGroup.Get(nil,key.(string),groupcache.AllocatingByteSliceSink(&data))
	object, err := bytesToInterface(data,this.T)
	if err != nil {
		return nil, err
	}
	return object,nil
}

func (this *GroupCacheKvStorage) Set(key interface{}, object interface{}) error {
	return nil
}

func (this *GroupCacheKvStorage) MultiGet(keys []interface{}) (map[interface{}]interface{}, error) {
	resultMap := make(map[interface{}]interface{})
	for _,key := range(keys){
		value,err:=this.Get(key)
		if err !=nil{
			continue
		}
		resultMap[key] = value
	}
	return resultMap,nil
}

func (this *GroupCacheKvStorage) MultiSet(objectMap map[interface{}]interface{}) error {
	return nil
}

func (this *GroupCacheKvStorage) Delete(key interface{}) error {
	return nil
}

func bytesToInterface(data []byte,T reflect.Type) (interface{}, error) {
	tStruct := reflect.New(T)
	dec := json.NewDecoder(bytes.NewBuffer(data))
	err := dec.Decode(tStruct.Interface())
	if err != nil {
		return err, nil
	}
	return reflect.Indirect(tStruct.Elem()).Interface(), nil
}

