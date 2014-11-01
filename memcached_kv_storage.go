package storage

import (
	"errors"
	"reflect"
	"strings"

	"github.com/dropbox/godropbox/memcache"
)

type MemcachedStorage struct {
	client            memcache.Client
	KeyPrefix         string
	DefaultExpireTime int
	encoding          Encoding
}

func NewMcStorage(client memcache.Client, keyPrefix string, defaultExpireTime int, encoding Encoding) MemcachedStorage {
	return MemcachedStorage{client, keyPrefix, defaultExpireTime, encoding}
}

func (this MemcachedStorage) Get(key Key) (interface{}, error) {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	item := this.client.Get(cacheKey)
	if item.Error() != nil || item.Status() != memcache.StatusNoError {
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
	if object == nil {
		return nil
	}
	if reflect.TypeOf(object).Kind() == reflect.Slice {
		s := reflect.ValueOf(object)
		if s.IsNil() {
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
	response := this.client.Set(&memcache.Item{Key: keyCache, Value: buf, Expiration: uint32(this.DefaultExpireTime)})
	return response.Error()
}

func (this MemcachedStorage) MultiGet(keys []Key) (map[Key]interface{}, error) {
	keyMap := make(map[Key]interface{})
	for _, key := range keys {
		keyMap[key] = nil
	}
	cacheKeys := make([]string, len(keyMap))
	i := 0
	for key, _ := range keyMap {
		cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
		if err != nil {
			return nil, err
		}
		cacheKeys[i] = cacheKey
		i = i + 1
	}
	itemMap := this.client.GetMulti(cacheKeys)
	result := make(map[Key]interface{})
	for k, item := range itemMap {
		if len(item.Value()) == 0 {
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
	items := make([]*memcache.Item, 0, len(objectMap))
	for k, v := range objectMap {
		buf, err := this.encoding.Marshal(v)
		if err != nil {
			return err
		}
		keyCache, err := BuildCacheKey(this.KeyPrefix, k)
		item := &memcache.Item{Key: keyCache, Value: buf, Expiration: uint32(this.DefaultExpireTime)}
		items = append(items, item)
	}
	responses := this.client.SetMulti(items)
	for _, response := range responses {
		if response.Error() != nil {
			return response.Error()
		}
	}
	return nil
}

func (this MemcachedStorage) Delete(key Key) error {
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return err
	}
	response := this.client.Delete(cacheKey)
	return response.Error()
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
