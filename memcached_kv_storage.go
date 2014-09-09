package storage

import (
	"reflect"
	"github.com/dropbox/godropbox/memcache"
	"strings"
	"errors"
)

type MemcachedStorage struct {
	client            memcache.Client
	KeyPrefix         string
	DefaultExpireTime int
	encoding          Encoding
}

func NewMcStorage(client memcache.Client, keyPrefix string, defaultExpireTime int,encoding Encoding) MemcachedStorage {
	return MemcachedStorage{client, keyPrefix, defaultExpireTime,encoding}
}

func (this MemcachedStorage) Get(key Key) (interface{}, error) {
	cacheKey, err :=BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	item := this.client.Get(cacheKey)
	if item.Error() != nil || item.Status()!=memcache.StatusNoError{
		return nil, item.Error()
	} else {
		object, err := this.encoding.Unmarshal(item.Value())
		if err != nil {
			return nil, err
		} else {
			return object, nil
		}
	}
}

func (this MemcachedStorage) Set(key Key, object interface{}) error {
	if object==nil{
		return nil
	}
	if reflect.TypeOf(object).Kind()==reflect.Slice{
		s := reflect.ValueOf(object)
		if(s.IsNil()){
			return nil
		}
	}
	buf, err := this.encoding.Marshal(object)
	if err != nil {
		return err
	}
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Set(&memcache.Item{Key: keyCache, Value: buf,Expiration:uint32(this.DefaultExpireTime)})
	return nil
}

func (this MemcachedStorage) MultiGet(keys []Key) (map[Key]interface{}, error) {
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
		if err != nil {
			return nil, err
		}
		cacheKeys[index] = cacheKey
	}
	itemMap:= this.client.GetMulti(cacheKeys)
	result := make(map[Key]interface{})
	for k, item := range itemMap {
		if (len(item.Value())==0) {
			continue
		}
		object, err := this.encoding.Unmarshal(item.Value())
		if err != nil {
			continue
		}
		result[GetRawKey(k)] = object
	}
	return result, nil
}

func (this MemcachedStorage) MultiSet(objectMap map[Key]interface{}) error {
	for k, v := range objectMap {
		if err := this.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (this MemcachedStorage) Delete(key Key) error {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Delete(cacheKey)
	return nil
}

func (this MemcachedStorage) FlushAll() {
	this.client.Flush(uint32(0))
}

func BuildCacheKey(keyPrefix string, key Key) (cacheKey string, err error) {
	if key == nil {
		return "", errors.New("key should not be nil")
	}
	return strings.Join([]string{keyPrefix, key.ToString()}, "_"), nil
}

func GetRawKey(key string) (rawKey String) {
	keys := strings.Split(key, "_")
	return String(keys[len(keys)-1])
}

func InitializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			InitializeStruct(ft.Type, f)
		default:
		}
	}
}

/*
type MemcachedStorage struct {
	client            *memcache.Client
	KeyPrefix         string
	DefaultExpireTime int
	encoding          Encoding
}

func NewMcStorage(serverUrls []string, keyPrefix string, defaultExpireTime int,encoding Encoding) *MemcachedStorage {
	client := memcache.New(serverUrls...)
	return &MemcachedStorage{client, keyPrefix, defaultExpireTime,encoding}
}

func (this MemcachedStorage) Get(key Key) (interface{}, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	item, err := this.client.Get(cacheKey)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		return nil, err
	} else {
		object, err := this.encoding.Unmarshal(item.Value)
		if err != nil {
			return nil, err
		} else {
			return object, nil
		}
	}
}

func (this MemcachedStorage) Set(key Key, object interface{}) error {
	if object==nil{
		return nil
	}
	if reflect.TypeOf(object).Kind()==reflect.Slice{
		s := reflect.ValueOf(object)
		if(s.IsNil()){
			return nil
		}
	}
	buf, err := this.encoding.Marshal(object)
	if err != nil {
		return err
	}
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Set(&memcache.Item{Key: keyCache, Value: buf,Expiration:int32(this.DefaultExpireTime)})
	return nil
}

func (this MemcachedStorage) MultiGet(keys []Key) (map[Key]interface{}, error) {
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
		if err != nil {
			return nil, err
		}
		cacheKeys[index] = cacheKey
	}
	itemMap, err := this.client.GetMulti(cacheKeys)
	if err != nil {
		return nil, err
	}
	result := make(map[Key]interface{})
	for k, item := range itemMap {
		object, err := this.encoding.Unmarshal(item.Value)
		if err != nil {
			continue
		}
		result[GetRawKey(k)] = object
	}
	return result, nil
}

func (this MemcachedStorage) MultiSet(objectMap map[Key]interface{}) error {
	for k, v := range objectMap {
		if err := this.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (this MemcachedStorage) Delete(key Key) error {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Delete(cacheKey)
	return nil
}

func (this MemcachedStorage) FlushAll() {
	this.client.FlushAll()
}

*/
