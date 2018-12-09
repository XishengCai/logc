package main

// context 包

import(
	"context"
	"fmt"
)
func main() {
	// gen在单独的goroutine中生成整数
	// 将他们呢送到返回的频道
	// gen的调用者需要取消一次上下文
	// 他们完成消耗生成的整数不泄露
	// 内部goroutine由gen开始
	gen := func(ctx context.Context) <- chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <- ctx.Done():
					return  // 返回不要泄漏goroutine
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	// 背景返回non-nil（非零），空的 Context。它从未被取消，没有值，也没有最后期限。它通常由主函数，
	// 初始化和测试使用，并作为传入请求的top-level Context （顶级上下文）。
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 当我们完成小号整数时取消

	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}