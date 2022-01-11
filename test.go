package main

import "fmt"

type t struct {
	a     int
	aa    int
	str   string
	slice []int
}

func main() {
	// tmp := make([]int, 1)
	// fmt.Println(tmp[1:1])
	obj := &t{a: 3, aa: 123123, str: "1asdasd3"}
	obj.slice = make([]int, 4)
	aa(*obj)
	return
}

func aa(s interface{}) {
	switch s.(type) {
	case t:
		fmt.Println("Type it t")
		tmp := s.(t)
		fmt.Printf("%+v", tmp)
	default:
		fmt.Println("Unknow")
	}
	return
}
