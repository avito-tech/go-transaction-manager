package txinfo_mutext_vs_atomic

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

type Propagation int

const (
	PropagationRequired Propagation = iota
	PropagationNested
)

type TxInfo struct {
	Propagation  Propagation
	IsNew        bool
	IsNested     bool
	NestingLevel int
}

/******** RWMutex version ********/

type rwState struct {
	mu sync.RWMutex
	tx TxInfo
}

func (st *rwState) read() {
	st.mu.RLock()
	_ = st.tx.Propagation
	_ = st.tx.IsNested
	_ = st.tx.NestingLevel
	st.mu.RUnlock()
}

func (st *rwState) readAndWrite() {
	st.mu.RLock()
	_ = st.tx.Propagation
	_ = st.tx.IsNested
	_ = st.tx.NestingLevel
	st.mu.RUnlock()

	st.mu.Lock()
	st.tx.NestingLevel++
	st.tx.IsNested = true
	st.tx.Propagation = PropagationNested
	st.mu.Unlock()
}

/******** atomic.Value version ********/

type atomicState struct {
	v atomic.Value // stores *TxInfo
}

func newAtomicState() *atomicState {
	st := &atomicState{}
	st.v.Store(&TxInfo{})
	return st
}

func (st *atomicState) read() {
	_ = st.v.Load().(*TxInfo)
}

func (st *atomicState) readAndWrite() {
	x := st.v.Load().(*TxInfo)

	_ = x.Propagation
	_ = x.IsNested
	_ = x.NestingLevel

	n := &TxInfo{
		Propagation:  PropagationNested,
		IsNew:        x.IsNew,
		IsNested:     true,
		NestingLevel: x.NestingLevel + 1,
	}
	st.v.Store(n)
}

/******** Benchmarks ********/

func BenchmarkRWMutex_TxInfo_3Read_1Write(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	st := &rwState{}
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.readAndWrite()
			st.read()
			st.read()
			st.readAndWrite()
		}
	})
}

func BenchmarkAtomicValue_TxInfo_3Read_1Write(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	st := newAtomicState()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.readAndWrite()
			st.read()
			st.read()
			st.readAndWrite()
		}
	})
}

/******** Mutex version ********/

type mutexState struct {
	mu sync.Mutex
	tx TxInfo
}

func (st *mutexState) read() {
	st.mu.Lock()
	_ = st.tx.Propagation
	_ = st.tx.IsNested
	_ = st.tx.NestingLevel
	st.mu.Unlock()
}

func (st *mutexState) readAndWrite() {
	st.mu.Lock()
	_ = st.tx.Propagation
	_ = st.tx.IsNested
	_ = st.tx.NestingLevel

	st.tx.NestingLevel++
	st.tx.IsNested = true
	st.tx.Propagation = PropagationNested
	st.mu.Unlock()
}

/******** Benchmarks ********/

func BenchmarkMutex_TxInfo_3Read_1Write(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	st := &mutexState{}
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			st.readAndWrite()
			st.read()
			st.read()
			st.readAndWrite()
		}
	})
}
