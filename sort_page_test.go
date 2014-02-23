package storage

import (
	"testing"
	"sort"
)

func TestPage(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(array)))
	result:=Page(IntReversedSlice(array),0,0,1,20)

	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}

	result=Page(IntReversedSlice(array),0,0,2,20)
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=180{
		t.Error("first one should be 180")
	}
	if result.(IntReversedSlice)[19]!=161{
		t.Error("first one should be 161")
	}

	result=Page(IntReversedSlice(array),20,0,1,20)
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}

	result=Page(IntReversedSlice(array),190,0,1,20)
	if result.(IntReversedSlice).Len()!=10{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[9]!=191{
		t.Error("first one should be 191")
	}

	result=Page(IntReversedSlice(array),0,190,1,20)
	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=189{
		t.Error("first one should be 189")
	}
	if result.(IntReversedSlice)[19]!=170{
		t.Error("first one should be 170")
	}

	result=Page(IntReversedSlice(array),140,153,1,20)
	if result.(IntReversedSlice).Len()!=12{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=152{
		t.Error("first one should be 152")
	}
	if result.(IntReversedSlice)[11]!=141{
		t.Error("first one should be 141")
	}

}

func TestSortPage(t *testing.T) {
	var array []int
	for i := 1; i <= 200; i++ {
		array=append(array,i)
	}
	result:=SortAndPage(IntReversedSlice(array),0,0,1,20)

	if result.(IntReversedSlice).Len()!=20{
		t.Error("len should be 20")
	}
	if result.(IntReversedSlice)[0]!=200{
		t.Error("first one should be 200")
	}
	if result.(IntReversedSlice)[19]!=181{
		t.Error("first one should be 181")
	}
}
