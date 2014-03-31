package storage

var MAXLEN=200

func (this *MemcachedStorage) Getlimit(key,sinceId,maxId interface{},page,count int)(interface{},error){
	obj,err:=this.Get(key)
	if err!=nil{
		return nil,err
	}
	return Page(obj.(Pagerable),sinceId,maxId,page,count),nil
}

func (this *MemcachedStorage) AddItem(key interface{},item interface{})error{
	obj,err:=this.Get(key)
	if err!=nil{
		return err
	}
	result:=obj.(Pagerable).AddItem(item,MAXLEN)
	return this.Set(key,result)
}

func (this *MemcachedStorage) DeleteItem(key interface{},item interface{})error{
	obj,err:=this.Get(key)
	if err!=nil{
		return err
	}
	result:=obj.(Pagerable).DeleteItem(item)
	return this.Set(key,result)
}
