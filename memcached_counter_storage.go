package storage

func (this *MemcachedKvStorage) Incr(key interface{},step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0,err
	}
	return this.client.Increment(keyCache,step)
}

func (this *MemcachedKvStorage) Decr(key interface{},step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0,err
	}
	return this.client.Decrement(keyCache,step)
}
