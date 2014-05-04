package storage

import (
	"reflect"
	"testing"
)

func TestIncrDecrRedis(t *testing.T) {
	jsonEncoding:=JsonEncoding{reflect.TypeOf(1)}
	redisStorage,_ := NewRedisStorage(":6379", "test", 0, jsonEncoding)
	redisStorage.Set(String("1"), 1)
	res, _ := redisStorage.Get(String("1"))
	defer redisStorage.Delete(String("1"))
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if res.(int)!=1{
		t.Error("value should be 1")
	}

	resIncr,_:=redisStorage.Incr(String("1"),1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=2{
		t.Error("value should be 2")
	}

	resIncr,_=redisStorage.Incr(String("1"),3)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=5{
		t.Error("value should be 5")
	}

	resDecr,_:=redisStorage.Decr(String("1"),1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=4{
		t.Error("value should be 4")
	}

	resDecr,_=redisStorage.Decr(String("1"),2)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=2{
		t.Error("value should be 2")
	}

	resDecr,err:=redisStorage.Decr(String("2"),2)
	if err!=nil{
		t.Error("err should be nil",err)
	}
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=0{
		t.Error("value should be 0")
	}

}
