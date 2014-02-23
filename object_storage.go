package storage

type Storage interface {
	Get(key interface{}) (interface{}, error)
	Set(key interface{}, object interface{}) error
	MultiGet(keys []interface{}) (map[interface{}]interface{}, error)
	MultiSet(map[interface{}]interface{}) error
	Delete(key interface{}) error
}

type CounterStorage interface{
	Storage
	Incr(key interface{},step uint64)(newValue uint64, err error)
	Decr(key interface{},step uint64)(newValue uint64, err error)
}

type ListStorage interface{
	Storage
	Getlimit(key ,sinceId ,maxId interface{},page ,count int)(interface{},error)
}

type StorageProxy struct {
	PreferedStorage Storage
	BackupStorage   Storage
}

func NewStorageProxy(prefered, backup Storage) *StorageProxy {
	return &StorageProxy{
		PreferedStorage: prefered,
		BackupStorage:   backup,
	}
}

func (this *StorageProxy) Get(key interface{}) (interface{}, error) {
	object, err := this.PreferedStorage.Get(key)
	if err != nil {
		return nil, err
	}
	if object == nil {
		object, err = this.BackupStorage.Get(key)
		if err != nil {
			return nil, err
		}
		if object != nil {
			this.PreferedStorage.Set(key, object)
		}
	}
	return object, nil
}

func (this *StorageProxy) Set(key interface{}, object interface{}) error {
	if object != nil {
		err := this.PreferedStorage.Set(key, object)
		if err != nil {
			return err
		}
		err = this.BackupStorage.Set(key, object)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *StorageProxy) MultiGet(keys []interface{}) (map[interface{}]interface{}, error) {
	resultMap, err := this.PreferedStorage.MultiGet(keys)
	if err != nil {
		return nil, err
	}
	missedKeyCount := 0
	for _, key := range keys {
		if _, find := resultMap[key]; !find {
			missedKeyCount++
		}
	}
	if missedKeyCount > 0 {
		missedKeys := make([]interface{}, missedKeyCount)
		i := 0
		for _, key := range keys {
			if _, find := resultMap[key]; !find {
				missedKeys[i] = key
				i++
			}
		}
		missedMap, err := this.BackupStorage.MultiGet(missedKeys)
		if err != nil {
			return nil, err
		}
		this.MultiSet(missedMap)
		for k, v := range missedMap {
			resultMap[k] = v
		}
	}
	return resultMap, nil
}

func (this *StorageProxy) MultiSet(objectMap map[interface{}]interface{}) error {
	err := this.PreferedStorage.MultiSet(objectMap)
	if err != nil {
		return err
	}
	err = this.BackupStorage.MultiSet(objectMap)
	if err != nil {
		return err
	}
	return nil
}

func (this *StorageProxy) Delete(key interface{}) error {
	err := this.BackupStorage.Delete(key)
	if err != nil {
		return err
	}
	err = this.PreferedStorage.Delete(key)
	if err != nil {
		return err
	}
	return nil
}

func (this *StorageProxy) Incr(key interface{},step uint64)(newValue uint64, err error){
	result,err:=this.PreferedStorage.(CounterStorage).Incr(key,step)
	if err!=nil{
		return result,err
	}
	result,err=this.BackupStorage.(CounterStorage).Incr(key,step)
	if err!=nil{
		return result,err
	}
	return result,err
}

func (this *StorageProxy) Decr(key interface{},step uint64)(newValue uint64, err error){
	result,err:=this.PreferedStorage.(CounterStorage).Decr(key,step)
	if err!=nil{
		return result,err
	}
	result,err=this.BackupStorage.(CounterStorage).Decr(key,step)
	if err!=nil{
		return result,err
	}
	return result,err
}
