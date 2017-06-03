package msgbus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueLen(t *testing.T) {
	q := Queue{}
	assert.Equal(t, q.Len(), 0)
}

func TestQueuePush(t *testing.T) {
	q := Queue{}
	q.Push(1)
	assert.Equal(t, q.Len(), 1)
}

func TestQueuePop(t *testing.T) {
	q := Queue{}
	q.Push(1)
	assert.Equal(t, q.Len(), 1)
	assert.Equal(t, q.Pop(), 1)
	assert.Equal(t, q.Len(), 0)
}

func TestQueuePeek(t *testing.T) {
	q := Queue{}
	q.Push(1)
	assert.Equal(t, q.Len(), 1)
	assert.Equal(t, q.Peek(), 1)
	assert.Equal(t, q.Len(), 1)
}

func TestQueue(t *testing.T) {
	q := Queue{}
	q.Push(1)
	q.Push(2)
	q.Push(3)
	assert.Equal(t, q.Len(), 3)
	assert.Equal(t, q.Peek(), 1)

	assert.Equal(t, q.Pop(), 1)
	assert.Equal(t, q.Len(), 2)
	assert.Equal(t, q.Peek(), 2)

	assert.Equal(t, q.Pop(), 2)
	assert.Equal(t, q.Len(), 1)
	assert.Equal(t, q.Peek(), 3)

	assert.Equal(t, q.Pop(), 3)
	assert.Equal(t, q.Len(), 0)
	assert.Equal(t, q.Peek(), nil)
}

func BenchmarkQueuePush(b *testing.B) {
	q := Queue{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkQueuePeekEmpty(b *testing.B) {
	q := Queue{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Peek()
	}
}

func BenchmarkQueuePeekNonEmpty(b *testing.B) {
	q := Queue{}
	q.Push(1)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Peek()
	}
}

func BenchmarkQueuePopEmpty(b *testing.B) {
	q := Queue{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}

func BenchmarkQueuePopNonEmpty(b *testing.B) {
	q := Queue{}
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Pop()
	}
}
