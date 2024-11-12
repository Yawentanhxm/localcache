package main

import (
	"fmt"
	"localcache/localcache"
	"time"
	"unsafe"
)

func main() {
	cache := localcache.MyCache{}
	cache.SetMaxMemory("16B")
	a := int8(0)
	fmt.Println(unsafe.Sizeof(a))
	cache.Set("11", "11", time.Second)
	cache.Set("22", "22", time.Second)
	cache.Print()
	//cache.Set("33","22",time.Second)
	cache.Get("11")
	cache.Get("11")
	cache.Get("11")
	cache.Get("22")
	cache.Print()
	cache.Set("44", "22", time.Second)
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
