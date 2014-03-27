package storage

import (
	"testing"
	"reflect"
	"sort"
)

func TestGetLimitRedis(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice:=IntReversedSlice(array)
	jsonEncoding:=JsonEncoding{reflect.TypeOf(slice)}
	redisStorage,_ := NewRedisStorage(":6379", "test", 0, jsonEncoding)
	redisListStorage:=RedisListStorage{redisStorage}
	redisListStorage.Set("1", slice)
	result, _ := redisListStorage.Getlimit("1",0,0,1,20)
	defer redisListStorage.Delete("1")
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}
	if string(result.([]interface{})[0].([]byte))!="200"{
		t.Error("first one should be 200")
	}
	if string(result.([]interface{})[19].([]byte))!="181"{
		t.Error("first one should be 181")
	}
}

func TestAddItemRedis(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice:=IntReversedSlice(array)
	jsonEncoding:=JsonEncoding{reflect.TypeOf(slice)}
	redisStorage,_ := NewRedisStorage(":6379", "test", 0, jsonEncoding)
	redisListStorage:=RedisListStorage{redisStorage}

	redisListStorage.Set("1", slice)
	result, _ := redisListStorage.Getlimit("1",0,0,1,20)
	defer redisListStorage.Delete("1")
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}
	if string(result.([]interface{})[0].([]byte))!="200"{
		t.Error("first one should be 200")
	}
	if string(result.([]interface{})[19].([]byte))!="181"{
		t.Error("first one should be 181")
	}

	redisListStorage.AddItem("1",201)
	result, _ = redisListStorage.Getlimit("1",0,0,1,20)
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}
	if string(result.([]interface{})[0].([]byte))!="201"{
		t.Error("first one should be 201")
	}
	if string(result.([]interface{})[19].([]byte))!="182"{
		t.Error("first one should be 182")
	}

	redisListStorage.DeleteItem("1",193)
	result, _ = redisListStorage.Getlimit("1",0,0,1,20)
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}
	if string(result.([]interface{})[0].([]byte))!="201"{
		t.Error("first one should be 201")
	}
	if string(result.([]interface{})[19].([]byte))!="181"{
		t.Error("first one should be 181")
	}

}
