package storage

import (
	"reflect"
	"testing"
)

func TestIncrDecr(t *testing.T) {
	jsonEncoding:=JsonEncoding{reflect.TypeOf(1)}
	mcStorage := NewMcStorage([]string{"localhost:12000"}, "test", 0, jsonEncoding)
	mcStorage.Set("1", 1)
	res, _ := mcStorage.Get("1")
	defer mcStorage.Delete("1")
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if res.(int)!=1{
		t.Error("value should be 1")
	}

	resIncr,_:=mcStorage.Incr("1",1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=2{
		t.Error("value should be 2")
	}

	resIncr,_=mcStorage.Incr("1",3)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resIncr!=5{
		t.Error("value should be 5")
	}

	resDecr,_:=mcStorage.Decr("1",1)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=4{
		t.Error("value should be 4")
	}

	resDecr,_=mcStorage.Decr("1",2)
	if reflect.TypeOf(res) != reflect.TypeOf(1) {
		t.Error("res type is not T")
	}
	if resDecr!=2{
		t.Error("value should be 2")
	}

	resDecr,err:=mcStorage.Decr("2",2)
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
