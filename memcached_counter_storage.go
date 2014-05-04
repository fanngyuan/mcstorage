package storage

import (
	"github.com/fanngyuan/gomemcache/memcache"
)

func (this MemcachedStorage) Incr(key interface{},step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return 0, nil
		}
		return 0,err
	}
	result,errcache:=this.client.Increment(keyCache,step)
	if errcache==memcache.ErrCacheMiss {
		return 0, nil
	}
	return result,errcache
}

func (this MemcachedStorage) Decr(key interface{},step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0,err
	}
	result,errcache:=this.client.Decrement(keyCache,step)
	if errcache==memcache.ErrCacheMiss {
		return 0, nil
	}
	return result,errcache
}
