package storage

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGetSetProxy(t *testing.T) {
	tt := T{1}
	mcStorage1 := NewMcStorage([]string{"localhost:12000"}, "test_1", 0, reflect.TypeOf(&tt))
	mcStorage2 := NewMcStorage([]string{"localhost:12000"}, "test_2", 0, reflect.TypeOf(&tt))
	storageProxy := NewStorageProxy(mcStorage1, mcStorage2)

	mcStorage1.Set("1", tt)
	res, _ := storageProxy.Get("1")
	defer storageProxy.Delete("1")
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
	res, _ = mcStorage2.Get("1")
	if res != nil {
		t.Error("there should be no 1 in mc2")
	}

	mcStorage2.Set("2", tt)
	res, _ = storageProxy.Get("2")
	defer storageProxy.Delete("2")
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	res, _ = mcStorage1.Get("2")
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
	mcStorage1 := NewMcStorage([]string{"localhost:12000"}, "test_1", 0, reflect.TypeOf(&tt))
	mcStorage2 := NewMcStorage([]string{"localhost:12000"}, "test_2", 0, reflect.TypeOf(&tt))
	storageProxy := &StorageProxy{mcStorage1, mcStorage2}

	valueMap := make(map[interface{}]interface{})
	keys := make([]interface{}, 10)
	for i := 0; i < 10; i++ {
		keys[i] = strconv.Itoa(i)
		valueMap[strconv.Itoa(i)] = T{i}
		defer storageProxy.Delete(strconv.Itoa(i))
	}
	mcStorage1.MultiSet(valueMap)
	res, _ := storageProxy.MultiGet(keys)
	for k, v := range res {
		if reflect.TypeOf(v) != reflect.TypeOf(tt) {
			t.Error("res type is not T")
		}
		kint, err := strconv.Atoi(k.(string))
		if err != nil {
			t.Error("key %s is not int ", k)
		}
		vT := v.(T)
		if kint != vT.A {
			t.Error("value should be %s,while it is %s", kint, vT.A)
		}
	}

	valueMap2 := make(map[interface{}]interface{})
	keys2 := make([]interface{}, 10)
	for i := 10; i < 20; i++ {
		keys2[i-10] = strconv.Itoa(i)
		valueMap[strconv.Itoa(i)] = T{i}
		defer storageProxy.Delete(strconv.Itoa(i))
	}

	mcStorage2.MultiSet(valueMap2)
	res, _ = storageProxy.MultiGet(keys2)
	for k, v := range res {
		if reflect.TypeOf(v) != reflect.TypeOf(tt) {
			t.Error("res type is not T")
		}
		kint, err := strconv.Atoi(k.(string))
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
		kint, err := strconv.Atoi(k.(string))
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
	mcStorage1 := NewMcStorage([]string{"localhost:12000"}, "test_1", 0, reflect.TypeOf(&tt))
	mcStorage2 := NewMcStorage([]string{"localhost:12000"}, "test_2", 0, reflect.TypeOf(&tt))
	storageProxy := &StorageProxy{mcStorage1, mcStorage2}

	mcStorage2.Set("2", tt)
	res, _ := storageProxy.Get("2")
	defer storageProxy.Delete("2")
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	res, _ = mcStorage1.Get("2")
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes = res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}

	storageProxy.Delete("2")
	res, _ = mcStorage1.Get("2")
	if res != nil {
		t.Error("res should be nil ,after delete")
	}

	res, _ = mcStorage2.Get("2")
	if res != nil {
		t.Error("res should be nil ,after delete")
	}

	res, _ = storageProxy.Get("2")
	if res != nil {
		t.Error("res should be nil ,after delete")
	}
}

func TestIncrDecrProxy(t *testing.T) {
	mcStorage1 := NewMcStorage([]string{"localhost:12000"}, "test_1", 0, reflect.TypeOf(1))
	mcStorage2 := NewMcStorage([]string{"localhost:12000"}, "test_2", 0, reflect.TypeOf(1))
	storageProxy := NewStorageProxy(mcStorage1, mcStorage2)

	mcStorage2.Set("1", 1)
	res, _ := storageProxy.Get("1")
	defer storageProxy.Delete("1")
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}

	resIncr,_:=storageProxy.Incr("1",1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=2{
		t.Error("value should be 2")
	}

	resDecr,_:=storageProxy.Decr("1",1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=1{
		t.Error("value should be 1")
	}

}
