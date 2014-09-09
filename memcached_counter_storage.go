package storage

func (this MemcachedStorage) Incr(key Key,step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0,err
	}
	result:=this.client.Increment(keyCache,step,0,0)
	if result.Error() != nil{
		return 0, result.Error()
	}
	return result.Count(),nil
}

func (this MemcachedStorage) Decr(key Key,step uint64)(newValue uint64, err error){
	keyCache, err := BuildCacheKey(this.KeyPrefix, key)
	if err != nil {
		return 0,err
	}
	result:=this.client.Decrement(keyCache,step,0,0)
	if result.Error() != nil{
		return 0, result.Error()
	}
	return result.Count(),nil
}



