package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"reflect"
	"strings"
)

type MemcachedKvStorage struct {
	client            *memcache.Client
	KeyPrefix         string
	DefaultExpireTime int
	T                 reflect.Type
}

func NewMcStorage(serverUrls []string, keyPrefix string, defaultExpireTime int, t reflect.Type) *MemcachedKvStorage {
	client := memcache.New(serverUrls...)
	return &MemcachedKvStorage{client, keyPrefix, defaultExpireTime, t}
}

func (this *MemcachedKvStorage) Get(key interface{}) (interface{}, error) {
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
		err, object := this.bytesToInterface(item)
		if err != nil {
			return nil, err
		} else {
			return object, nil
		}
	}
}

func (this *MemcachedKvStorage) bytesToInterface(item *memcache.Item) (error, interface{}) {
	tStruct := reflect.New(this.T)
	dec := json.NewDecoder(bytes.NewBuffer(item.Value))
	err := dec.Decode(tStruct.Interface())
	if err != nil {
		return err, nil
	}
	return nil, tStruct.Elem().Interface()
}

func (this *MemcachedKvStorage) Set(key interface{}, object interface{}) error {
	buf, err := json.Marshal(object)
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

func (this *MemcachedKvStorage) MultiGet(keys []interface{}) (map[interface{}]interface{}, error) {
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
		err, object := this.bytesToInterface(item)
		if err != nil {
			continue
		}
		result[GetRawKey(k)] = object
	}
	return result, nil
}

func (this *MemcachedKvStorage) MultiSet(objectMap map[interface{}]interface{}) error {
	for k, v := range objectMap {
		if err := this.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func (this *MemcachedKvStorage) Delete(key interface{}) error {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	this.client.Delete(cacheKey)
	return nil
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
