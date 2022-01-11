package raft

import (
	"fmt"
	"log"
)

// Debugging
const (
	Debug   = false
	Verbose = true
)

var Name = []string{"Leader", "Follower", "Candidate"}

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

func (rf *Raft) dump() {
	if Debug {
		if !Verbose {
			fmt.Println("---------------------------------------------------------")
			fmt.Printf("raft: %d, serverState: %s\n", rf.me, Name[rf.serverState])
			fmt.Printf("raft's nextIndex: ")
			fmt.Println(rf.volStateOnLdr.nextIndex)
		} else {
			fmt.Printf("%#v\n", rf)
		}
	}
}

func dumpArgs(i interface{}) {
	if Debug {
		switch i.(type) {
		case AppendEntriesArgs:
			obj := i.(AppendEntriesArgs)
			fmt.Printf("%#v\n", obj)
		case AppendEntriesReply:
			obj := i.(AppendEntriesReply)
			fmt.Printf("%#v\n", obj)
		case RequestVoteArgs:
			obj := i.(RequestVoteArgs)
			fmt.Printf("%#v\n", obj)
		case RequestVoteReply:
			obj := i.(RequestVoteReply)
			fmt.Printf("%#v\n", obj)
		default:
			fmt.Printf("Unknown type\n")
		}
	}
}

func AssertEqual(a interface{}, b interface{}, msg string) bool {
	if Debug {
		if a != b {
			fmt.Printf(msg)
		}
	}
	return a == b
}

func AssertNotEqual(a interface{}, b interface{}, msg string) bool {
	if Debug {
		if a == b {
			fmt.Printf(msg)
		}
	}
	return a != b
}

func AssertBigger(a int, b int, msg string) bool {
	if Debug {
		if a <= b {
			fmt.Printf(msg)
		}
	}
	return a > b
}
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
