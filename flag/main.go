package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if false {
		showOsArgs()
	}
	showFlag()

}

func showFlag() {
	var (
		name string
		age  int
	)
	flag.StringVar(&name, "name", "张三", "姓名")
	flag.IntVar(&age, "age", 18, "年龄")

	flag.Parse()

	fmt.Printf("亲爱的%s 招商银行祝您%d岁生日快乐! \n", name, age)

}

func showOsArgs() {
	fmt.Println(os.Args)
	if len(os.Args) > 0 {
		for i, v := range os.Args {
			fmt.Printf("args[%d]=%v\n", i, v)
		}
	}
}
