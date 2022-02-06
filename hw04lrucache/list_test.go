package hw04lrucache

import (
	"reflect"
	"testing"
)

type lTest struct {
	wantLen   int
	wantFront interface{}
	wantBack  interface{}
}

func equalIntElemsList(t *testing.T, l List, want []int) {
	elems := make([]int, 0, l.Len())
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}

	if !reflect.DeepEqual(elems, want) {
		t.Errorf("elems: %v; want: %v", elems, want)
	}
}

func equalLenFrontBackList(t *testing.T, l List, test lTest) {
	var frontList, backList interface{}
	lenList := l.Len()

	if l.Front() != nil {
		frontList = l.Front().Value
	}
	if l.Back() != nil {
		backList = l.Back().Value
	}

	if lenList != test.wantLen {
		t.Errorf("list.Len() = %d; want: %d", lenList, test.wantLen)
	}
	if frontList != test.wantFront {
		t.Errorf("list.Front() = %v; want: %v", frontList, test.wantFront)
	}
	if backList != test.wantBack {
		t.Errorf("list.Back() = %v; want: %v", backList, test.wantBack)
	}
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		equalLenFrontBackList(t, l, lTest{0, nil, nil})
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		equalLenFrontBackList(t, l, lTest{3, 10, 30})

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]

		equalLenFrontBackList(t, l, lTest{2, 10, 30})

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		equalLenFrontBackList(t, l, lTest{7, 80, 70})

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		want := []int{70, 80, 60, 40, 10, 30, 50}
		equalIntElemsList(t, l, want)
	})

	t.Run("list shift", func(t *testing.T) {
		sl := []int{80, 60, 40, 10, 30, 50, 70}
		l := NewList()

		for _, v := range sl {
			l.PushBack(v)
		} // [80, 60, 40, 10, 30, 50, 70]

		for i := 0; i < len(sl); i++ {
			l.MoveToFront(l.Front())
		}

		want := []int{80, 60, 40, 10, 30, 50, 70}
		equalIntElemsList(t, l, want)
	})
}
