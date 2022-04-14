package hw04lrucache

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
	sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache { //nolint:ireturn
	return &lruCache{capacity, NewList(), make(map[Key]*listItem, capacity), sync.Mutex{}}
}

// Set is adds a value to the cache.
func (cache *lruCache) Set(key Key, value interface{}) bool {
	cache.Lock()
	defer cache.Unlock()

	itemList, isNotCacheMiss := cache.items[key]
	if isNotCacheMiss {
		itemCache, _ := itemList.Value.(*cacheItem)
		itemCache.value = value
		cache.queue.MoveToFront(itemList)
	} else {
		newItemCache := &cacheItem{key: key, value: value}
		newItemList := cache.queue.PushFront(newItemCache)
		cache.items[key] = newItemList
		if cache.queue.Len() > cache.capacity {
			lastItemQueue := cache.queue.Back()
			lastItemCacheQueue, _ := lastItemQueue.Value.(*cacheItem)
			cache.queue.Remove(lastItemQueue)
			delete(cache.items, lastItemCacheQueue.key)
		}
	}

	return isNotCacheMiss
}

// Get is gets a value from the cache.
func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.Lock()
	defer cache.Unlock()

	itemList, isNotCacheMiss := cache.items[key]
	if isNotCacheMiss {
		itemCache, _ := itemList.Value.(*cacheItem)
		cache.queue.MoveToFront(itemList)

		return itemCache.value, isNotCacheMiss
	}

	return nil, isNotCacheMiss
}

// Clear is clears the cache.
func (cache *lruCache) Clear() {
	cache.Lock()
	defer cache.Unlock()

	for key := range cache.items {
		cache.queue.Remove(cache.items[key])
		delete(cache.items, key)
	}
}
