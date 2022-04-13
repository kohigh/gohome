package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	_, exists := l.items[key]

	if exists {
		l.items[key].Value = cacheItem{key, value}
		l.queue.MoveToFront(l.items[key])

		return true
	}

	l.items[key] = l.queue.PushFront(cacheItem{key, value})
	if l.capacity < l.queue.Len() {
		delete(l.items, l.queue.Back().Value.(cacheItem).key)
		l.queue.Remove(l.queue.Back())
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	val, exists := l.items[key]

	if exists {
		l.queue.MoveToFront(val)
		return val.Value.(cacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
