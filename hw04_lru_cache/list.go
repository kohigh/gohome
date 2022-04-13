package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type listItemState int

const (
	used listItemState = iota
	removed
)

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
	state listItemState
}

type listState int

const (
	idle listState = iota
	single
	multiple
)

const (
	incr string = "incr"
	decr        = "decr"
)

type list struct {
	size  int
	state listState
	front *ListItem
	back  *ListItem
}

func (l list) Len() int {
	return l.size
}

func (l list) Front() *ListItem {
	return l.front
}

func (l list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	oldFront := l.front
	l.front = &ListItem{Value: v, Next: oldFront}

	l.updState(incr)

	if l.state == multiple {
		oldFront.Prev = l.front
	}

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	oldBack := l.back
	l.back = &ListItem{Value: v, Prev: oldBack}

	l.updState(incr)

	if l.state == multiple {
		oldBack.Next = l.back
	}

	return l.back
}

func (l *list) updState(action string) {
	switch action {
	case incr:
		if l.size++; l.size > 1 {
			l.state = multiple
		}
	case decr:
		if l.size--; l.size <= 1 {
			l.state = single
		}
	}

	if l.front == nil {
		l.front = l.back
	}

	if l.back == nil {
		l.back = l.front
	}
}

func (l *list) Remove(i *ListItem) {
	if i.state == removed {
		return
	}

	l.updState(decr)

	if l.front == i {
		l.front = i.Next
	}

	if l.back == i {
		l.back = i.Prev
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = nil
	i.state = removed
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}

func NewList() List {
	return new(list)
}
