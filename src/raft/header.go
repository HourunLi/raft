package raft

import (
	"math/rand"
	"sync"
	"time"

	"6.824/labrpc"
)

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in part 2D you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh, but set CommandValid to false for these
// other uses.
//

type ServerState int32

const (
	Leader ServerState = iota
	Follower
	Candidate
)

const CHANSIZE = 1 << 6

const (
	AppendEntryInterval = time.Duration(100 * time.Millisecond) //send heartbeat per 0.1s
)

//the time limit for each selection round is 0.5-1s
func electionTimeout() time.Duration {
	return time.Duration(800+rand.Intn(800)) * time.Millisecond
}

// *usage send a ApplyMsg to the service when commiting a new log
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int

	// For 2D:
	SnapshotValid bool
	Snapshot      []byte
	SnapshotTerm  int
	SnapshotIndex int
}

/**********************************************************.
|                                                          |
|    example RequestVote RPC arguments/reply structure.    |
|      field names must start with capital letters         |
|                                                          |
`**********************************************************/

type AppendEntriesArgs struct {
	Term         int        // leader's term
	LeaderId     int        // so follower can redirect clients
	PrevLogIndex int        // index of log entry immediately preceding new ones
	PrevLogTerm  int        // term of prevLogIndex entry
	Entries      []LogEntry //log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit int        //leader's commit index
}

type AppendEntriesReply struct {
	Term         int  // currentTerm, for leader to update itself
	Succeed      bool // true if follower contained entry matching prevLogIndex and prevLogTerm
	NextTryIndex int  // the next attempt index for appending
}

type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term         int // candidate’s term
	CandidateId  int // candidate requesting vote
	LastLogIndex int // index of candidate’s last log entry
	LastLogTerm  int // term of candidate’s last log entry
}

type RequestVoteReply struct {
	// Your data here (2A).
	Term        int  // currentTerm, for candidate to update itself
	VoteGranted bool // true means candidate received vote
}

type InstallSnapshotArgs struct {
	Term              int
	LeaderId          int
	LastIncludedIndex int
	LastIncludedTerm  int
	Data              []byte
}

type InstallSnapshotReply struct {
	Term int
}

type LogEntry struct {
	LogIndex int
	LogTerm  int
	Command  interface{}
}

type PersistState struct {
	CurrentTerm int        // latest term server has seen
	VoteFor     int        // candidateId that received vote in current term
	VoteCnt     int        // the number of votes got
	Logs        []LogEntry // log entries;
}

type VolatileStateOnServers struct {
	commitIndex int // index of highest log entry known to be committed
	applyIndex  int // index of highest log entry known to be committed
}

type VolatileStateOnLeaders struct {
	nextIndex  []int // index of highest log entry known to be committed
	matchIndex []int // for each server, index of highest log entry known to be replicated on server
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu            sync.Mutex             // Lock to protect shared access to this peer's state
	peers         []*labrpc.ClientEnd    // RPC end points of all peers
	persister     *Persister             // Object to hold this peer's persisted state
	me            int                    // this peer's index into peers[]
	dead          int32                  // set by Kill()
	leaderId      int                    // leader's id
	serverState   ServerState            // server's identity state
	perState      PersistState           // PersistState
	volStateOnSer VolatileStateOnServers // Volatile State On Servers
	volStateOnLdr VolatileStateOnLeaders // Volatile State On Leaders
	applyCh       chan ApplyMsg          // apply to client
	heartBeat     chan struct{}          // heart beat channel
	winElection   chan struct{}          // win election channel
	grantVote     chan struct{}          // vote channel
	shutDown      chan struct{}          // shutDown channel
}
