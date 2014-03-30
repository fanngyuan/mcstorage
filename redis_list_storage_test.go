package storage

import (
	"testing"
	"sort"
)

func TestGetLimitRedis(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	slice:=IntReversedSlice(array)
	redisListStorage,_ := NewRedisListStorage(":6379", "test", 0, DecodeIntReversedSlice)
	redisListStorage.Set("1", slice)
	result, _ := redisListStorage.Getlimit("1",0,0,1,20)
	defer redisListStorage.Delete("1")
	if len(result.(IntReversedSlice))!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
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
	redisListStorage,_ := NewRedisListStorage(":6379", "test", 0, DecodeIntReversedSlice)

	redisListStorage.Set("1", slice)
	result, _ := redisListStorage.Getlimit("1",0,0,1,20)
	defer redisListStorage.Delete("1")
	if len(result.(IntReversedSlice))!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}

	result, _ = redisListStorage.Getlimit("1",0,200,1,20)
	if len(result.(IntReversedSlice))!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=199{
		t.Error("first one should be 199")
	}
	if result.(IntReversedSlice)[19]!=180{
		t.Error("first one should be 180")
	}

	redisListStorage.AddItem("1",201)
	result, _ = redisListStorage.Getlimit("1",0,0,1,20)
	if len(result.(IntReversedSlice))!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=201{
		t.Error("first one should be 201")
	}
	if result.(IntReversedSlice)[19]!=182{
		t.Error("first one should be 182")
	}

	redisListStorage.DeleteItem("1",193)
	result, _ = redisListStorage.Getlimit("1",0,0,1,20)
	if len(result.(IntReversedSlice))!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=201{
		t.Error("first one should be 201")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}

}
