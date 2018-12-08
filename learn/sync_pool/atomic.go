package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// sync/atomic 库提供了原子库操作的支持,原子操作直接有底层CPU硬件支持，因而一般要比基于操作系统API的锁方式高效。

func main() {
	var sum uint32 = 100
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			//sum += 1
			atomic.AddUint32(&sum, 1)
		}()
	}
	wg.Wait()
	fmt.Println(sum)
}

// 函数原型：
// atomic.AddUint32(addr *uint32, delta uint32) uint32
// atomic.AddUint64(addr *uint64, delta uint64) uint64
// atomic.AddUintptr(addr *uintptr, delta uintptr) uintptr

