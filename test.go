package main

import "fmt"

func main() {
	tmp := make([]int, 0)
	for i := range tmp {
		fmt.Println(i)
	}
	return
}
