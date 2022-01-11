package raft

/******************************************************.
|                The structure of logs                 |
|------------------------------------------------------|
|    LogIndex      0         1         2       ...     |
|    LogTerm       0         ?         ?       ...     |
|    Command    padding    command   command   ...     |
`******************************************************/

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	//	"bytes"

	"sync/atomic"
	"time"

	//	"6.824/labgob"
	"6.824/labrpc"
)

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool
	// Your code here (2A).
	term = rf.perState.currentTerm
	isleader = rf.serverState == Leader
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
// A service wants to switch to snapshot.  Only do so if Raft hasn't
// have more recent info since it communicate the snapshot on applyCh.
//
func (rf *Raft) CondInstallSnapshot(lastIncludedTerm int, lastIncludedIndex int, snapshot []byte) bool {

	// Your code here (2D).

	return true
}

// the service says it has created a snapshot that has
// all info up to and including index. this means the
// service no longer needs the log through (and including)
// that index. Raft should now trim its log as much as possible.
func (rf *Raft) Snapshot(index int, snapshot []byte) {
	// Your code here (2D).
}

func (rf *Raft) applyLogs() {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	for i := rf.volStateOnSer.applyIndex + 1; i <= rf.volStateOnSer.commitIndex; i++ {
		msg := ApplyMsg{CommandValid: true, CommandIndex: i, Command: rf.perState.logs[i].Command}
		if i > 0 {
			rf.applyCh <- msg
		}
		rf.volStateOnSer.applyIndex = i
	}
	return
}

//
//judge whether RequestVoteArgs is at least as new as rf's state
//
func (rf *Raft) argsIsUpToDate(args *RequestVoteArgs) bool {
	lastTerm := rf.getLastLogTerm()
	lastIndex := rf.getLastLogIndex()
	if lastTerm < args.Term ||
		(lastTerm == args.Term && lastIndex <= args.LastLogIndex) {
		return true
	}
	return false
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term < rf.perState.currentTerm {
		reply.Term = rf.perState.currentTerm
		reply.VoteGranted = false
		return
	}
	rf.perState.currentTerm = args.Term
	reply.Term = args.Term
	rf.serverState = Follower
	if (rf.perState.voteFor == -1 || rf.perState.voteFor == args.CandidateId) && rf.argsIsUpToDate(args) {
		rf.perState.voteFor = args.CandidateId
		reply.VoteGranted = true
		rf.grantVote <- struct{}{}
		return
	}
	reply.VoteGranted = false
	return
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if !AssertEqual(rf.serverState, Candidate, "rf's state should be Candidate\n") {
		return ok
	}
	if !AssertEqual(rf.perState.currentTerm, args.Term, "rf's term should equals to args\n") {
		return ok
	}
	if ok {
		if reply.Term > rf.perState.currentTerm {
			rf.serverState = Follower
			rf.perState.currentTerm = reply.Term
			rf.perState.voteFor = -1
			return ok
		}
		// if reply term and rf term is the same
		// it means the server gets the vote request in the present vote round
		if reply.Term == rf.perState.currentTerm {
			if reply.VoteGranted {
				rf.perState.voteCnt++
				if rf.perState.voteCnt > len(rf.peers)/2 {
					rf.serverState = Leader
					// rf.volStateOnLdr = VolatileStateOnLeaders{}
					rf.volStateOnLdr.matchIndex = make([]int, len(rf.peers))
					rf.volStateOnLdr.nextIndex = make([]int, len(rf.peers))
					tmp := len(rf.perState.logs)
					DPrintf("%d\n", tmp)
					for i := range rf.peers {
						rf.volStateOnLdr.nextIndex[i] = tmp
					}
					rf.winElection <- struct{}{}
				}
			}
		}
	}
	return ok
}

func (rf *Raft) sendRequestVotes() {
	if !AssertEqual(rf.serverState, Candidate, "rf's state should be Candidate\n") {
		return
	}
	rf.mu.Lock()
	defer rf.mu.Unlock()
	requestVoteArgs := &RequestVoteArgs{}
	requestVoteArgs.Term = rf.perState.currentTerm
	requestVoteArgs.CandidateId = rf.me
	requestVoteArgs.LastLogIndex = rf.getLastLogIndex()
	requestVoteArgs.LastLogTerm = rf.getLastLogTerm()
	for server := range rf.peers {
		if server != rf.me {
			go rf.sendRequestVote(server, requestVoteArgs, &RequestVoteReply{})
		}
	}
}

// candidate: get the heartbeat, switch identity to follower and fresh term etc.
func (rf *Raft) appendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	// The append request is stale due to sluggish or diordered network
	if args.Term < rf.perState.currentTerm {
		reply.Term = rf.perState.currentTerm
		reply.Succeed = false
		reply.NextTryIndex = rf.getNextTryIndex()
		return
	}

	// args.Term == rf.perState.currentTerm is possible

	// maybe rf's state is candidate, which requires to switch state to follower
	if args.Term > rf.perState.currentTerm {
		rf.serverState = Follower
		rf.perState.currentTerm = args.Term
		rf.perState.voteFor = -1
	}

	rf.heartBeat <- struct{}{}
	reply.Term = rf.perState.currentTerm

	// if rf's(follower) log is shorter than leader's
	// then set leader's try index straightly to rf(follower) logs' tail
	if args.PrevLogIndex > rf.getLastLogIndex() {
		reply.NextTryIndex = rf.getNextTryIndex()
		reply.Succeed = false
		return
	}

	// bypass the inconsistent entries of the same term every time
	if args.PrevLogIndex > 0 && args.PrevLogTerm != rf.perState.logs[args.PrevLogIndex].LogTerm {
		term := rf.perState.logs[args.PrevLogIndex].LogTerm
		for i := args.PrevLogIndex; i >= 0; i-- {
			if rf.perState.logs[i].LogTerm == term {
				continue
			}
			reply.NextTryIndex = i + 1
			break
		}
	} else {
		rf.perState.logs = rf.perState.logs[:args.PrevLogIndex+1]
		rf.perState.logs = append(rf.perState.logs, args.Entries...)
		reply.Succeed = true
		reply.NextTryIndex = rf.getNextTryIndex()

		if rf.volStateOnSer.commitIndex < args.LeaderCommit {
			rf.volStateOnSer.commitIndex = Min(args.LeaderCommit, rf.getLastLogIndex())
			go rf.applyLogs()
		}
	}
	return
}

func (rf *Raft) sendAppendEntriesSignal(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
	ok := rf.peers[server].Call("Raft.appendEntries", args, reply)
	if !AssertEqual(rf.serverState, Leader, "rf's state should be Leader\n") {
		return ok
	}
	if !AssertEqual(rf.perState.currentTerm, args.Term, "rf's current term should be equal to args term\n") {
		return ok
	}
	if !ok {
		DPrintf("sendAppendEntriesSignal Failed\n")
	}
	if reply.Term > rf.perState.currentTerm {
		rf.serverState = Follower
		rf.perState.voteFor = -1
		rf.perState.currentTerm = reply.Term
		return ok
	}

	if reply.Succeed {
		lens := len(args.Entries)
		if lens > 0 {
			rf.volStateOnLdr.nextIndex[server] = args.Entries[lens-1].LogIndex + 1
			rf.volStateOnLdr.matchIndex[server] = args.Entries[lens-1].LogIndex
		}
	} else {
		rf.volStateOnLdr.nextIndex[server] = reply.NextTryIndex
	}

	for cmtIndex := rf.getLastLogIndex(); cmtIndex > rf.volStateOnSer.commitIndex &&
		rf.perState.logs[cmtIndex].LogTerm == rf.perState.currentTerm; cmtIndex-- {
		cnt := 1
		for i := range rf.peers {
			if i != rf.me && rf.volStateOnLdr.matchIndex[i] >= cmtIndex {
				cnt++
			}
		}
		if cnt > len(rf.peers)/2 {
			rf.volStateOnSer.commitIndex = cmtIndex
			go rf.applyLogs()
			break
		}
	}
	return ok
}

//
// Leader send append entries to the followers
//
func (rf *Raft) sendAppendEntriesSignals() {
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if !AssertEqual(rf.serverState, Leader, "rf's state should be Leader\n") {
		return
	}
	for server := range rf.peers {
		if server == rf.me {
			continue
		}
		args := &AppendEntriesArgs{}
		args.Term = rf.perState.currentTerm
		args.LeaderId = rf.me
		args.PrevLogIndex = rf.volStateOnLdr.nextIndex[server] - 1
		args.PrevLogTerm = rf.perState.logs[args.PrevLogIndex].LogTerm
		args.Entries = rf.perState.logs[rf.volStateOnLdr.nextIndex[server]:]
		args.LeaderCommit = rf.volStateOnSer.commitIndex
		go rf.sendAppendEntriesSignal(server, args, &AppendEntriesReply{})
	}
	return
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	// Your code here (2B).
	rf.mu.Lock()
	term := rf.perState.currentTerm
	index := len(rf.perState.logs)
	isLeader := rf.serverState == Leader
	rf.mu.Unlock()

	if isLeader {
		go func(rf *Raft) {
			rf.mu.Lock()
			defer rf.mu.Unlock()
			rf.perState.logs = append(rf.perState.logs, LogEntry{LogIndex: index, LogTerm: term, Command: command})
			return
		}(rf)
	}

	return index, term, isLeader
}

//
// the tester doesn't halt goroutines created by Raft after each test,
// but it does call the Kill() method. your code can use killed() to
// check whether Kill() has been called. the use of atomic avoids the
// need for a lock.
//
// the issue is that long-running goroutines use memory and may chew
// up CPU time, perhaps causing later tests to fail and generating
// confusing debug output. any goroutine with a long-running loop
// should call killed() to check whether it should stop.
//
func (rf *Raft) Kill() {
	atomic.StoreInt32(&rf.dead, 1)
	// Your code here, if desired.
}

func (rf *Raft) killed() bool {
	z := atomic.LoadInt32(&rf.dead)
	return z == 1
}

// The ticker go routine starts a new election if this peer hasn't received
// heartsbeats recently.
func (rf *Raft) ticker() {
	for rf.killed() == false {
		// Your code here to check if a leader election should
		// be started and to randomize sleeping time using
		// time.Sleep().
		switch rf.serverState {
		case Leader:
			go rf.sendAppendEntriesSignals()
			time.Sleep(AppendEntryInterval)
		case Follower:
			select {
			case <-rf.grantVote:
			case <-rf.heartBeat:
			case <-time.After(electionTimeout()):
				rf.mu.Lock() // necessary to lock?
				rf.serverState = Candidate
				rf.mu.Unlock()
			}
		case Candidate:
			rf.mu.Lock()
			rf.perState.voteFor = rf.me
			rf.perState.currentTerm++
			rf.perState.voteCnt = 1
			rf.mu.Unlock()
			go rf.sendRequestVotes() // send request votes

			select {
			case <-rf.heartBeat:
				rf.mu.Lock() // necessary to lock?
				rf.serverState = Follower
				rf.mu.Unlock()
			case <-rf.winElection:
			case <-time.After(electionTimeout()):
			}
		}
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.leaderId = -1
	rf.serverState = Follower
	// for persistent state
	rf.perState.currentTerm = 0
	rf.perState.voteFor = -1
	rf.perState.logs = append(rf.perState.logs, LogEntry{LogIndex: 0, LogTerm: 0, Command: nil}) // padding
	rf.perState.voteCnt = 0
	// for volatile state on server
	rf.volStateOnSer.applyIndex = 0
	rf.volStateOnSer.commitIndex = 0

	// volatile state on leader initialization until the rf server becomes leader
	// ! attention: not initialized
	// rf.volStateOnLdr.matchIndex = make([]int, len(peers))
	// rf.volStateOnLdr.nextIndex = make([]int, len(peers))

	// for channel
	rf.applyCh = applyCh
	rf.heartBeat = make(chan struct{}, CHANSIZE)
	rf.winElection = make(chan struct{}, CHANSIZE)
	rf.grantVote = make(chan struct{}, CHANSIZE)

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// go apply
	// start ticker goroutine to start elections
	go rf.ticker()

	return rf
}
