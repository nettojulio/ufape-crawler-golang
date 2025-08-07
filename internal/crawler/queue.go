package crawler

import "github.com/emirpasic/gods/queues/linkedlistqueue"

type QueueItem struct {
	URL   string
	Depth uint
}

type Queue struct {
	q *linkedlistqueue.Queue
}

func NewQueue() *Queue {
	return &Queue{q: linkedlistqueue.New()}
}

func (q *Queue) Enqueue(item QueueItem) {
	q.q.Enqueue(item)
}

func (q *Queue) Dequeue() QueueItem {
	itemRaw, _ := q.q.Dequeue()
	return itemRaw.(QueueItem)
}

func (q *Queue) Empty() bool {
	return q.q.Empty()
}
