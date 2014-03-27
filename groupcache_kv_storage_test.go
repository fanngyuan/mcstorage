package storage

import (
	"reflect"
	"testing"
	"github.com/golang/groupcache"
	"encoding/json"
)

func TestGetSetGC(t *testing.T) {
	tt := T{1}

	jsonEncoding:=JsonEncoding{reflect.TypeOf(&tt)}
	mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)
	mcStorage.Set("1", tt)
	res, _ := mcStorage.Get("1")
	defer mcStorage.Delete("1")
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	var groupcache = groupcache.NewGroup("SlowDBCache", 64<<20, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result,err := mcStorage.Get(key)
			if err!=nil{
				return nil
			}
			bytes,err:=json.Marshal(result)
			if err!=nil{
				return nil
			}
			dest.SetBytes(bytes)
			return nil
		}))
	gcStorage := &GroupCacheKvStorage{groupcache,0,jsonEncoding}
	res,_=gcStorage.Get("1")
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	mcStorage.Delete("1")
	res,_=gcStorage.Get("1")
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
}
