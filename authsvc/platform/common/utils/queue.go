package utils

import (
	"container/list"
	"sync"
)

type Queue struct {
	l *list.List
	lock sync.Mutex
}

// Init initializes the queue data structure.
// A queue must be initialized before it can be used.
// O(1)
func (q *Queue) Init() {
	q.l = list.New()
}

// Push enqueues an element to the queue.
// O(1)
func (q *Queue) Push(v interface{}) {
	q.lock.Lock()
	defer q.lock.Unlock()
	
	q.l.PushFront(v)
}

// Pop dequeues an element from the queue.
// O(1)
func (q *Queue) Pop() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	
	if q.l.Len() == 0 {
		return nil
	}

	v := q.l.Back()
	return q.l.Remove(v)
}

// Len returns the number of elements in the queue.
// O(1)
func (q *Queue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	
	return q.l.Len()
}

// IsEmpty returns true the queue has no elements.
// O(1)
func (q *Queue) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	
	return q.l.Len() == 0
}

func (q *Queue) Peek() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()
	
	if q.l.Len() == 0 {
		return nil
	}

	v := q.l.Back()
	
	return v.Value;
}