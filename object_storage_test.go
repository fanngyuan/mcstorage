package storage

import (
	"reflect"
	"strconv"
	"testing"
	"github.com/dropbox/godropbox/memcache"
)

func TestGetSetProxy(t *testing.T) {
	tt := T{1}
	jsonEncoding:=JsonEncoding{reflect.TypeOf(&tt)}

	client:=memcache.NewMockClient()
	mcStorage1 := NewMcStorage(client, "test_1", 0, jsonEncoding)
	mcStorage2 := NewMcStorage(client, "test_2", 0, jsonEncoding)

	storageProxy := NewStorageProxy(mcStorage1, mcStorage2)

	mcStorage1.Set(String("1"), tt)
	res, _ := storageProxy.Get(String("1"))
	defer storageProxy.Delete(String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
	res, _ = mcStorage2.Get(String("1"))
	if res != nil {
		t.Error("there should be no 1 in mc2")
	}

	mcStorage2.Set(String("2"), tt)
	res, _ = storageProxy.Get(String("2"))
	defer storageProxy.Delete(String("2"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	res, _ = mcStorage1.Get(String("2"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
}

func TestMultiGetSetProxy(t *testing.T) {
	tt := T{1}

	jsonEncoding:=JsonEncoding{reflect.TypeOf(&tt)}
	client:=memcache.NewMockClient()
	mcStorage1 := NewMcStorage(client, "test_1", 0, jsonEncoding)
	mcStorage2 := NewMcStorage(client, "test_2", 0, jsonEncoding)

	storageProxy := &StorageProxy{mcStorage1, mcStorage2}

	valueMap := make(map[Key]interface{})
	keys := make([]Key, 10)
	for i := 0; i < 10; i++ {
		keys[i] = String(strconv.Itoa(i))
		valueMap[String(strconv.Itoa(i))] = T{i}
		defer storageProxy.Delete(String(strconv.Itoa(i)))
	}
	mcStorage1.MultiSet(valueMap)
	res, _ := storageProxy.MultiGet(keys)
	for k, v := range res {
		if reflect.TypeOf(v) != reflect.TypeOf(tt) {
			t.Error("res type is not T")
		}
		kint, err := strconv.Atoi(k.ToString())
		if err != nil {
			t.Error("key %s is not int ", k)
		}
		vT := v.(T)
		if kint != vT.A {
			t.Error("value should be %s,while it is %s", kint, vT.A)
		}
	}

	valueMap2 := make(map[Key]interface{})
	keys2 := make([]Key, 10)
	for i := 10; i < 20; i++ {
		keys2[i-10] = String(strconv.Itoa(i))
		valueMap2[String(strconv.Itoa(i))] = T{i}
		defer storageProxy.Delete(String(strconv.Itoa(i)))
	}

	mcStorage2.MultiSet(valueMap2)
	res, _ = storageProxy.MultiGet(keys2)
	for k, v := range res {
		if reflect.TypeOf(v) != reflect.TypeOf(tt) {
			t.Error("res type is not T")
		}
		kint, err := strconv.Atoi(k.ToString())
		if err != nil {
			t.Error("key %s is not int ", k)
		}
		vT := v.(T)
		if kint != vT.A {
			t.Error("value should be %s,while it is %s", kint, vT.A)
		}
	}

	res, _ = mcStorage1.MultiGet(keys2)
	for k, v := range res {
		if reflect.TypeOf(v) != reflect.TypeOf(tt) {
			t.Error("res type is not T")
		}
		kint, err := strconv.Atoi(k.ToString())
		if err != nil {
			t.Error("key %s is not int ", k)
		}
		vT := v.(T)
		if kint != vT.A {
			t.Error("value should be %s,while it is %s", kint, vT.A)
		}
	}

}

func TestDeleteProxy(t *testing.T) {
	tt := T{1}

	jsonEncoding:=JsonEncoding{reflect.TypeOf(&tt)}
	client:=memcache.NewMockClient()
	mcStorage1 := NewMcStorage(client, "test_1", 0, jsonEncoding)
	mcStorage2 := NewMcStorage(client, "test_2", 0, jsonEncoding)
	storageProxy := &StorageProxy{mcStorage1, mcStorage2}

	mcStorage2.Set(String("2"), tt)
	res, _ := storageProxy.Get(String("2"))
	defer storageProxy.Delete(String("2"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	res, _ = mcStorage1.Get(String("2"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	storageProxy.Delete(String("2"))
	res, _ = mcStorage1.Get(String("2"))
	if res != nil {
		t.Error("res should be nil ,after delete")
	}

	res, _ = mcStorage2.Get(String("2"))
	if res != nil {
		t.Error("res should be nil ,after delete")
	}

	res, _ = storageProxy.Get(String("2"))
	if res != nil {
		t.Error("res should be nil ,after delete")
	}
}

/**
func TestIncrDecrProxy(t *testing.T) {
	jsonEncoding:=JsonEncoding{reflect.TypeOf(1)}

	client:=memcache.NewMockClient()
	mcStorage1 := NewMcStorage(client, "test_1", 0, jsonEncoding)
	mcStorage2 := NewMcStorage(client, "test_2", 0, jsonEncoding)

	storageProxy := NewStorageProxy(mcStorage1, mcStorage2)

	mcStorage2.Set(String("1"), 1)
	res, _ := storageProxy.Get(String("1"))
	defer storageProxy.Delete(String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}

	resIncr,_:=storageProxy.Incr(String("1"),1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=2{
		t.Error("value should be 2")
	}

	resDecr,_:=storageProxy.Decr(String("1"),1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=1{
		t.Error("value should be 1")
	}

}*/
