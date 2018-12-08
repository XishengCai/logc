package main

import (
	"log"
	"runtime"
)

func main(){
	test()
}

func test(){
	test2()
}

func test2(){
	// 前2层函数信息
	pc, file, line, ok := runtime.Caller(2)

	// 是否成功
	log.Println(ok)

	// 函数指针
	log.Println(pc)

	// 所属文件
	log.Println(file)

	// 所属行
	log.Println(line)

	// 获取函数信息
	f := runtime.FuncForPC(pc)

	// 函数名
	log.Println(f.Name())

	// 当前函数信息
	pc, file, line, ok = runtime.Caller(0)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f = runtime.FuncForPC(pc)
	log.Println(f.Name())

	// 上一层函数信息
	pc, file, line, ok = runtime.Caller(1)
	log.Println(pc)
	log.Println(file)
	log.Println(line)
	log.Println(ok)
	f = runtime.FuncForPC(pc)
	log.Println(f.Name())
}