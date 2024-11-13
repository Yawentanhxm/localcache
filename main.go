package main

import (
	"fmt"
	"localcache/localcache"
	"time"
	"unsafe"
)

func main() {
	cache := localcache.NewCache()
	cache.SetMaxMemory("16KB")
	a := int8(0)
	fmt.Println(unsafe.Sizeof(a))
	cache.Set("11", "11", time.Second)
	cache.Set("22", "22", 10*time.Second)
	cache.Set("33", "33", 3*time.Second)
	time.Sleep(2 * time.Second)
	cache.Print()
	time.Sleep(5 * time.Second)
	//cache.Set("33","22",time.Second)
	cache.Get("11")
	cache.Get("11")
	cache.Get("11")
	cache.Get("22")
	cache.Print()
	cache.Set("44", "22", time.Second)
	cache.Print()
	fmt.Println(cache.Keys())
	cache.Flush()
	fmt.Println(cache.Keys())
	cache.Set("11", "11", time.Second)
	cache.Set("22", "22", time.Second)
	fmt.Println(cache.Keys())
	cache.Print()
	//for ;;{
	//	//temp := pq.Pop()
	//	temp := heap.Pop(&pq)
	//	fmt.Println(temp)
	//	if temp == nil {
	//		return
	//	}
	//	//fmt.Println(heap.Pop(&pq))
	//}

}
