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
	return time.Duration(500 + rand.Intn(500))
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

type AppendEntries struct {
	term         int        // leader's term
	leaderId     int        // so follower can redirect clients
	prevLogIndex int        // index of log entry immediately preceding new ones
	prevLogTerm  int        // term of prevLogIndex entry
	entries      []LogEntry //log entries to store (empty for heartbeat; may send more than one for efficiency)
	leaderCommit int        //leader's commit index
}

type AppendEntriesReply struct {
	term    int  // currentTerm, for leader to update itself
	succeed bool // true if follower contained entry matching prevLogIndex and prevLogTerm
}

type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	term         int // candidate’s term
	candidateId  int // candidate requesting vote
	lastLogIndex int // index of candidate’s last log entry
	lastLogTerm  int // term of candidate’s last log entry
}

type RequestVoteReply struct {
	// Your data here (2A).
	term        int  // currentTerm, for candidate to update itself
	voteGranted bool // true means candidate received vote
}

type LogEntry struct {
	LogIndex int
	LogTerm  int
	command  interface{}
}

type PersistState struct {
	currentTerm int        // latest term server has seen
	voteFor     int        // candidateId that received vote in current term
	voteCnt     int        // the number of votes got
	logs        []LogEntry // log entries;
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
	shutDown      chan struct{}          // shutDown channel
}
