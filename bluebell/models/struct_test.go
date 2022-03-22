package models

import (
	"fmt"
	"testing"
	"unsafe"
)

type s1 struct {
	a int8
	b string
	c int8
}

type s2 struct {
	a int8
	c int8
	b string
}

//Go 内存对齐
//结构体定义时 尽量把相同类型的字段放在一起
func TestStruct(t *testing.T) {
	//unsafe.Sizeof(v1) : 32
	v1 := s1{
		a: 1,
		b: "qimi",
		c: 2,
	}

	//unsafe.Sizeof(v2) : 24
	v2 := s2{
		a: 1,
		c: 2,
		b: "qimi",
	}

	fmt.Println("unsafe.Sizeof(v1) : ", unsafe.Sizeof(v1))
	fmt.Println("unsafe.Sizeof(v2) :", unsafe.Sizeof(v2))
}
