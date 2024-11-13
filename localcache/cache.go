package localcache

import (
	"container/heap"
	"encoding/json"
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
	ch        <-chan time.Time
}

func NewCache() MyCache {
	cache := MyCache{}
	return cache
}

func (c *MyCache) timerExpr(key string, ch <-chan time.Time) {
	for {
		select {
		case <-ch:
			fmt.Printf("%s到定时时间\n", key)
			c.Del(key)
			return
		}
	}
}

type st struct {
	data interface{}         // 保存数据
	st   *priority_list.Item //保存key和使用计数
	//timer *time.Timer
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
			default:
				a = 0
			}
			break
		}
	}
	c.maxMemory = uintptr(a)
	heap.Init(&c.heap)
	return true
}

// todo 支持过期时间
// 加锁
func (c *MyCache) Set(key string, val interface{}, expire time.Duration) bool {

	item := &priority_list.Item{
		Value:    key,
		Priority: 1,
	}
	st := &st{
		data: val,
		st:   item,
	}
	timers := time.NewTimer(expire)
	c.data.Store(key, st)
	go c.timerExpr(key, timers.C)
	heap.Push(&c.heap, item)
	// todo 覆盖写的时候，Size不是累加
	bytes, _ := json.Marshal(st)
	newSize := unsafe.Sizeof(bytes)
	for c.curMemory+newSize > c.maxMemory {
		pop := heap.Pop(&c.heap)
		if pop == nil {
			return false
		}
		bytes, _ := json.Marshal(pop)
		c.curMemory -= unsafe.Sizeof(bytes)
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
	if value, ok := c.data.Load(key); ok {
		// 内存计数修改
		c.curMemory -= unsafe.Sizeof(value)
		// map中删除
		c.data.Delete(key)
		// 优先队列中删除
		heap.Remove(&c.heap, value.(*st).st.Index)
		return true
	}
	return false
}

func (c *MyCache) Exists(key string) bool {
	_, ok := c.data.Load(key)
	return ok
}

func (c *MyCache) Flush() bool {
	c.curMemory = 0
	// 直接指向一个新的
	c.heap = priority_list.PriorityQueue{}
	c.data = sync.Map{}
	return true
}

func (c *MyCache) Keys() int64 {
	return int64(c.heap.Len())
}

func (c *MyCache) Print() {
	for _, v := range c.heap {
		fmt.Println(v)
	}
}
