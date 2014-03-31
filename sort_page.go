package storage

import (
	"sort"
	"reflect"
)

type Pagerable interface{
	Find(key interface{})int
	Cut(start,end int)interface{}
	sort.Interface
	AddItem(value interface{},maxLen int)interface{}
	DeleteItem(value interface{})interface{}
}

func SortAndPage(pagerable Pagerable,sinceId,maxId interface{},page,count int)interface{}{
	sort.Sort(pagerable)
	return Page(pagerable,sinceId,maxId,page,count)
}

func Page(pagerable Pagerable,sinceId,maxId interface{},page,count int)interface{}{
	var start,end int
	i:=pagerable.Find(sinceId)
	if i>0{
		end=i-1
	}else{
		end=pagerable.Len()-1
	}
	i=pagerable.Find(maxId)
	start=i+1

	start=start+(page-1)*count

	if start>end && end!=0{
		return nil
	}
	var countReal int
	if (end-start+1)>count{
		countReal=count
	}else{
		countReal=(end-start+1)
	}
	end=start+countReal
	return pagerable.Cut(start,end)
}

type IntReversedSlice []int

//for reversed slice for [5,4,3,2,1]
func (this IntReversedSlice)Find(key interface{}) int{
	var k int
	if reflect.ValueOf(key).Type().Kind()==reflect.Uint{
		k=int(reflect.ValueOf(key).Uint())
		if k==0{
			return -1
		}
	}
	if reflect.ValueOf(key).Type().Kind()==reflect.Int{
		k=int(reflect.ValueOf(key).Int())
		if k==0{
			return -1
		}
	}
	i := sort.Search(len(this), func(i int) bool { return this[i] <= int(k) })
	if i < len(this) && this[i] == key {
		return i
	} else {
		return -1
	}
}

func (this IntReversedSlice)Cut(start,end int)interface{}{
	return this[start:end]
}

func (p IntReversedSlice) Len() int           { return len(p) }
func (p IntReversedSlice) Less(i, j int) bool { return p[i] > p[j] }
func (p IntReversedSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IntReversedSlice) AddItem(value interface{},maxLen int)interface{}{
	if p.Len()+1>maxLen{
		result:=make([]int,maxLen)
		result[0]=value.(int)
		copy(result[1:],p[:maxLen-1])
		return IntReversedSlice(result)
	}else{
		result:=make([]int,p.Len()+1)
		result[0]=value.(int)
		copy(result[1:],p)
		return IntReversedSlice(result)
	}
}

func (p IntReversedSlice) DeleteItem(value interface{})interface{}{
	index:=p.Find(value)
	if index==-1{
		return nil
	}
	result:=make([]int,p.Len()-1)
	copy(result,p[0:index])
	copy(result[index:],p[index+1:])
	return IntReversedSlice(result)
}
