package storage

import (
	"reflect"
	"strconv"
	"testing"
)

func TestGetSetRedis(t *testing.T) {
	tt := T{1}

	jsonEncoding:=JsonEncoding{reflect.TypeOf(&tt)}
	redisStorage,_ := NewRedisStorage(":6379", "test", 0, jsonEncoding)
	redisStorage.Set(String("1"), tt)
	res, _ := redisStorage.Get(String("1"))
	defer redisStorage.Delete(String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
}

func TestMultiGetSetRedis(t *testing.T) {
	tt := T{1}
	jsonEncoding:=JsonEncoding{reflect.TypeOf(&tt)}
	redisStorage,_ := NewRedisStorage(":6379", "test", 0, jsonEncoding)
	valueMap := make(map[Key]interface{})
	keys := make([]Key, 10)
	for i := 0; i < 10; i++ {
		keys[i] = String(strconv.Itoa(i))
		valueMap[String(strconv.Itoa(i))] = T{i}
		defer redisStorage.Delete(String(strconv.Itoa(i)))
	}
	redisStorage.MultiSet(valueMap)
	res, _ := redisStorage.MultiGet(keys)
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

func TestGetSetDeleteRedis(t *testing.T) {
	tt := T{1}
	jsonEncoding:=JsonEncoding{reflect.TypeOf(&tt)}
	redisStorage,_ := NewRedisStorage(":6379", "test", 0, jsonEncoding)
	redisStorage.Set(String("1"), tt)
	res, _ := redisStorage.Get(String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(tt) {
		t.Error("res type is not T")
	}
	ttRes := res.(T)
	if ttRes.A != tt.A {
		t.Error("res A field is not equals tt field")
	}
	redisStorage.Delete(String("1"))
	res, _ = redisStorage.Get(String("1"))

	if res!=nil{
		t.Error("res should be nil ,after delete")
	}
}
