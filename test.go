package main

import "fmt"

type Log struct {
	index   int
	term    int
	command interface{}
}

func main() {
	tmp := make([]Log, 0, 20)
	tmp = append(tmp, Log{index: 0, term: 0})
	tmp = append(tmp, Log{index: 1, term: 0})
	tmp = append(tmp, Log{index: 2, term: 1})
	tmp = append(tmp, Log{index: 3, term: 1})
	tmp2 := make([]Log, 0, 20)
	tmp2 = append(tmp2, tmp[0])
	fmt.Printf("%p, %p", &tmp[0], &tmp2[0])
	return
}
func shrinkLogs(logs []int) []int {
	// fmt.Printf("srclogs: %v", logs)
	if cap(logs) > len(logs)*2 {
		newLogs := make([]int, len(logs))
		copy(newLogs, logs)
		// fmt.Printf("newloogs: %v", newLogs)
		return newLogs
	}
	return logs
}
