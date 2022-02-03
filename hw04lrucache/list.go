package hw04lrucache

type List interface {
	Len() int
	Front() *listItem
	Back() *listItem
	PushFront(v interface{}) *listItem
	PushBack(v interface{}) *listItem
	Remove(i *listItem)
	MoveToFront(i *listItem)
}

// listItem is a list item.
type listItem struct {
	Value interface{} // value of list item.
	Next  *listItem   // next item in list.
	Prev  *listItem   // previous item in list.
}

// list is list :).
type list struct {
	len  int
	head *listItem
	tail *listItem
}

// NewList creates an empty list.
func NewList() List { //nolint:ireturn
	return &list{len: 0, head: nil, tail: nil}
}

// Len returns length of list.
func (l *list) Len() int {
	return l.len
}

// Front returns first item in list.
func (l *list) Front() *listItem {
	return l.head
}

// Back returns last item in list.
func (l *list) Back() *listItem {
	return l.tail
}

// PushFront adds item to head of list.
func (l *list) PushFront(v interface{}) *listItem {
	newItem := listItem{Value: v, Next: l.head, Prev: nil}

	if l.tail == nil {
		l.tail = &newItem
	} else {
		l.head.Prev = &newItem
	}
	l.head = &newItem

	l.len++

	return l.head
}

// PushBack adds item to tail of list.
func (l *list) PushBack(v interface{}) *listItem {
	newItem := listItem{Value: v, Next: nil, Prev: l.tail}

	if l.head == nil {
		l.head = &newItem
	} else {
		l.tail.Next = &newItem
	}
	l.tail = &newItem

	l.len++

	return l.tail
}

// Remove removes a list item.
func (l *list) Remove(i *listItem) {
	if l.head == i {
		l.head = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if l.tail == i {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

// MoveToFront moves item to front of list.
func (l *list) MoveToFront(i *listItem) {
	if l.head == i {
		return
	}

	if l.tail == i {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i.Prev.Next = i.Next

	i.Prev = nil
	i.Next = l.head

	l.head.Prev = i

	l.head = i
}
