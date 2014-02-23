package storage

import (
	"github.com/bradfitz/gomemcache/memcache"
)

func (this *MemcachedStorage) Getlimit(key,sinceId,maxId interface{},page,count int)(interface{},error){
	cacheKey, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return nil, err
	}
	item, err := this.client.Get(cacheKey)
	var obj interface{}
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		return nil, err
	} else {
		obj, err = this.bytesToInterface(item)
		if err != nil {
			return nil, err
		}
	}
	return Page(obj.(Pagerable),sinceId,maxId,page,count),nil
}
