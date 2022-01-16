package raft

import (
	"fmt"
	"log"
)

// Debugging
const (
	Debug   = false
	Verbose = true
	Info    = false
)

const StowageFactor = 2

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
			fmt.Printf("raft's nextIndex: %d", rf.volStateOnLdr.nextIndex)
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

func (rf *Raft) getFirstLog() LogEntry {
	return rf.perState.Logs[0]
}

func (rf *Raft) getLastLog() LogEntry {
	lens := len(rf.perState.Logs)
	return rf.perState.Logs[lens-1]
}

func (rf *Raft) getLastLogTerm() int {
	return rf.getLastLog().LogTerm
}

func (rf *Raft) getLastLogIndex() int {
	return rf.getLastLog().LogIndex
}

func (rf *Raft) getNextTryIndex() int {
	return rf.getLastLogIndex() + 1
}

func (rf *Raft) getFirstLogTerm() int {
	return rf.getFirstLog().LogTerm
}

func (rf *Raft) getFirstLogIndex() int {
	return rf.getFirstLog().LogIndex
}

//judge whether RequestVoteArgs is at least as new as rf's state
func (rf *Raft) argsIsUpToDate(args *RequestVoteArgs) bool {
	lastTerm, lastIndex := rf.getLastLogTerm(), rf.getLastLogIndex()
	DPrintf("argsIsUpToDate: lastTerm: %d, lastIndex: %d", lastTerm, lastIndex)
	if lastTerm < args.LastLogTerm ||
		(lastTerm == args.LastLogTerm && lastIndex <= args.LastLogIndex) {
		return true
	}
	return false
}

// If stowage factor is lower than 0.5, then shrink the log
func shrinkLogs(logs []LogEntry) []LogEntry {
	if cap(logs) > len(logs)*StowageFactor {
		newLogs := make([]LogEntry, len(logs))
		copy(newLogs, logs)
		return newLogs
	}
	return logs
}
