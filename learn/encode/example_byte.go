package main

import (
	"bytes"
	"fmt"
)

func main() {
	s1 := []byte("hello")       // 申请一个slice
	buff := bytes.NewBuffer(s1) // new 一个缓冲器buff, 里面存着hello这5个byte
	s2 := []byte("world")
	buff.Write(s2)
	fmt.Println(buff.String())

	s3 := make([]byte, 3)
	buff.Read(s3)
	fmt.Println("----",buff.String())
	fmt.Println(string(s3))

	buff.Read(s3)  // s3 == "low"
	fmt.Println(buff.String())
	fmt.Println(string(s3))
}