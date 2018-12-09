package main

import "unsafe"

type notifyList struct {
	wait uint32
	notify uint32
	lock uintptr
	head unsafe.Pointer
	tail unsafe.Pointer
}