package raft

import "log"

// Debugging
const Debug = true

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

func AssertEqual(a interface{}, b interface{}, msg string) bool {
	if Debug {
		if a != b {
			log.Printf(msg)
		}
	}
	return a == b
}

func (rf *Raft) getLastLogTerm() int {
	lens := len(rf.perState.logs)
	return rf.perState.logs[lens-1].LogTerm
}

func (rf *Raft) getLastLogIndex() int {
	lens := len(rf.perState.logs)
	return rf.perState.logs[lens-1].LogIndex
}

func (rf *Raft) getNextTryIndex() int {
	return len(rf.perState.logs)
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
