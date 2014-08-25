package storage

import (
	"testing"
)

func TestExistsMock(t *testing.T) {
	client:=NewMockClient()

	client.Set("aaa",[]byte("bbb"))
	result:=client.Exists("aaa")
	if result!=true{
		t.Errorf("result should be true")
	}

	result=client.Exists("bbb")
	if result!=false{
		t.Errorf("result should be false")
	}

	client.ClearAll()
	client.Set("aaa",[]byte("bbb"))
	v,err:=client.Get("aaa")
    if err!=nil {
		t.Errorf("%s\r\n",err.Error())
    }

	if string(v)!="bbb"{
		t.Errorf("result should be bbb")
	}

	client.ClearAll()
	client.Incr("aaa",1)
	v,err=client.Get("aaa")
	if string(v)!="1"{
		t.Errorf("result should be 1")
	}
	client.ClearAll()

	kvMap:=make(map[string][]byte)
	kvMap["1"]=[]byte("123")
	kvMap["2"]=[]byte("456")
	client.MultiSet(kvMap)

	v,err=client.Get("1")
	if string(v)!="123"{
		t.Errorf("result should be 123")
	}
	client.ClearAll()

}
