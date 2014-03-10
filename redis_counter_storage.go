package storage

import (
)

func (this RedisStorage) Incr(key interface{},step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0,err
	}
	result,errcache:=this.client.Incr(keyCache,step)
	return uint64(result),errcache
}

func (this RedisStorage) Decr(key interface{},step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0,err
	}
	result,errcache:=this.client.Decr(keyCache,step)
	if result<0{
		return 0,err
	}
	return uint64(result),errcache
}
