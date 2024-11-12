package localcache

import (
	"container/heap"
	"fmt"
	"localcache/priority_list"
	"sync"
	"time"
	"unsafe"
)

// 并发安全的本地缓存

type Cache interface {
	SetMaxMemory(size string) bool
	Set(key string, val interface{}, expire time.Duration) bool
	Get(key string) (interface{}, bool)
	Del(key string) bool
	Exists(key string) bool
	Flush() bool
	Keys() int64
}

type MyCache struct {
	data      sync.Map
	curMemory uintptr //Byte
	maxMemory uintptr //Byte
	heap      priority_list.PriorityQueue
}

type st struct {
	data interface{}
	st   *priority_list.Item
}

func (c *MyCache) SetMaxMemory(size string) bool {
	a := int64(0)
	for k, v := range size {
		if v <= '9' && v >= '0' {
			a = a*10 + int64(v-'0')
		} else {
			switch size[k:] {
			case "B":
			case "KB":
				a *= 1024
			case "MB":
				a *= 1024 * 1024
			case "GB":
				a *= 1024 * 1024 * 1024
			case "TB":
				a *= 1024 * 1024 * 1024 * 1024
			}
		}
	}
	c.maxMemory = uintptr(a)
	heap.Init(&c.heap)
	return true
}

// todo 支持过期时间
func (c *MyCache) Set(key string, val interface{}, expire time.Duration) bool {

	item := &priority_list.Item{
		Value:    key,
		Priority: 1,
	}
	st := &st{
		data: val,
		st:   item,
	}
	c.data.Store(key, st)
	heap.Push(&c.heap, item)
	newSize := unsafe.Sizeof(st)
	for c.curMemory+newSize > c.maxMemory {
		pop := heap.Pop(&c.heap)
		if pop == nil {
			return false
		}
		c.curMemory -= unsafe.Sizeof(pop)
	}
	c.curMemory += newSize
	return true
}

func (c *MyCache) Get(key string) (interface{}, bool) {
	if value, ok := c.data.Load(key); ok {
		value.(*st).st.Priority += 1
		heap.Fix(&c.heap, value.(*st).st.Index)
		return value, true
	}
	return nil, false
}

func (c *MyCache) Del(key string) bool {
	panic("implement me")
}

func (c *MyCache) Exists(key string) bool {
	panic("implement me")
}

func (c *MyCache) Flush() bool {
	panic("implement me")
}

func (c *MyCache) Keys() int64 {
	panic("implement me")
}
func (c *MyCache) Print() {
	for _, v := range c.heap {
		fmt.Println(v)
	}
}
