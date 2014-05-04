package storage

import (
	"errors"
	"github.com/fanngyuan/gomemcache/memcache"
	"reflect"
	"strings"
)

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

func (this MemcachedStorage) Get(key interface{}) (interface{}, error) {
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

func (this MemcachedStorage) Set(key interface{}, object interface{}) error {
	buf, err := this.encoding.Marshal(object)
	if err != nil {
		return err
	}
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Set(&memcache.Item{Key: keyCache, Value: buf})
	return nil
}

func (this MemcachedStorage) MultiGet(keys []interface{}) (map[interface{}]interface{}, error) {
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
	result := make(map[interface{}]interface{})
	for k, item := range itemMap {
		object, err := this.encoding.Unmarshal(item.Value)
		if err != nil {
			continue
		}
		result[GetRawKey(k)] = object
	}
	return result, nil
}

func (this MemcachedStorage) MultiSet(objectMap map[interface{}]interface{}) error {
	for k, v := range objectMap {
		if err := this.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (this MemcachedStorage) Delete(key interface{}) error {
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

func BuildCacheKey(keyPrefix interface{}, key interface{}) (cacheKey string, err error) {
	if key == nil {
		return "", errors.New("key should not be nil")
	}
	if reflect.TypeOf(key).Kind() != reflect.String {
		return "", errors.New("key should be string")
	}

	return strings.Join([]string{keyPrefix.(string), key.(string)}, "_"), nil
}

func GetRawKey(key string) (rawKey string) {
	keys := strings.Split(key, "_")
	return keys[len(keys)-1]
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
