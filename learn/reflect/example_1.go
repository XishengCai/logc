package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID		int
	Name	string
	Age 	int
}

func (u User) Hello() {
	fmt.Println("Hello cai")
}

// 定义一个反射函数，参数为任意类型
func Info(o interface{}) {
	// 使用反射类型获取o的Type, 一个包含多个方法的interface
	t := reflect.TypeOf(o)

	fmt.Printf("type: %s\n", t.Name())

	// 使用反射类型获取o的value， 一个空的结构体
	v := reflect.ValueOf(o)
	fmt.Println("Fields:")

	// t.NumField()打印结构体o的字段个数(ID, Name, Age共三个)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		fmt.Printf("%6s:%v = %v\n", f.Name, f.Type, val)
	}


	// 使用t.NumMethod()获取所有结构体类型的方法
	// 结构Type的方法NumMethod() int
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Printf("%6s:%v\n", m.Name, m.Type)
	}
}

func main() {
	u := User{1, "caixisheng", 25}
	Info(u)
	u.Hello()

}

//https://www.jianshu.com/p/42c19f88df6c