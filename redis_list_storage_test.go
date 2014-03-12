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
	redisStorage,_ := NewRedisStorage(":6379", "test", 0, reflect.TypeOf(&slice))
	redisListStorage:=RedisListStorage{redisStorage}
	redisListStorage.Set("1", slice)
	result, _ := redisListStorage.Getlimit("1",0,0,1,20)
	defer redisListStorage.Delete("1")
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}
}

func TestAddItemRedis(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice:=IntReversedSlice(array)
	redisStorage ,_:= NewRedisStorage(":6379", "test", 0, reflect.TypeOf(&slice))
	redisListStorage:=RedisListStorage{redisStorage}

	redisListStorage.Set("1", slice)
	result, _ := redisListStorage.Getlimit("1",0,0,1,20)
	defer redisListStorage.Delete("1")
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}

	redisListStorage.AddItem("1",201)
	result, _ = redisListStorage.Getlimit("1",0,0,1,20)
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}

	redisListStorage.DeleteItem("1",193)
	result, _ = redisListStorage.Getlimit("1",0,0,1,20)
	if len(result.([]interface{}))!=20{
		t.Error("len should be 20")
	}

}
