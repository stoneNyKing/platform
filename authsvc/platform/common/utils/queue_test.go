package utils

import (
	"testing"
	"fmt"
)

const iterations = 1024

func TestPushPop(t *testing.T) {
	q := new(Queue)
	q.Init()

	for i := 0; i < iterations; i++ {
		q.Push(i)
		fmt.Printf("push = %v\n",i)
	}

	for i := 0; i < iterations; i++ {
		testPop(t, q, i)
	}
}

func TestLen(t *testing.T) {
	q := new(Queue)
	q.Init()

	for i := 0; i < iterations; i++ {
		q.Push(i)
	}

	if l := q.Len(); l != iterations {
		t.Errorf("Queue length was expected to be %v, but is %v", iterations, l)
	}
	
	fmt.Printf("(TestLen)Queue len= %d\n",q.Len())
	
	q.Pop()
	if l := q.Len(); l != iterations-1 {
		t.Errorf("Queue length was expected to be %v, but is %v", iterations-1, l)
	}
}

func TestIsEmpty(t *testing.T) {
	q := new(Queue)
	q.Init()

	if q.IsEmpty() != true {
		fmt.Printf("Queue should be empty.\n")
	}else{
		fmt.Printf("(TestIsEmpty)Queue is empty.\n")
	}

	q.Push(1)
	fmt.Printf("(TestIsEmpty)Queue len= %d\n",q.Len())

	if q.IsEmpty() == false {
		fmt.Printf("Queue should not be empty\n")
	}

}

func testPop(t *testing.T, q *Queue, e interface{}) {
	if v := q.Pop(); v != e {
		t.Errorf("Popping expected %v, got %v", e, v)
	}
	
	fmt.Printf("popup = %v\n",e)
}

func BenchmarkPush(b *testing.B) {
	q := new(Queue)
	q.Init()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkPop(b *testing.B) {
	q := new(Queue)
	q.Init()

	for i := 0; i < b.N; i++ {
		q.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
