package server

import (
	"sync"
)

// Broadcaster acts as hub for broadcating snips to multiple recipients
type Broadcaster struct {
	sync.Mutex // guard recipients
	wg         sync.WaitGroup
	in         <-chan struct{}
	recipients []chan struct{}
	done       chan bool
}

// NewBroadcaster creates a Broadcaster that implements
// a hub and spoke message replication pattern
func NewBroadcaster(in <-chan struct{}) *Broadcaster {
	return &Broadcaster{
		in:         in,
		recipients: make([]chan struct{}, 0),
		done:       make(chan bool),
	}
}

// Run executes the broadcaster
func (b *Broadcaster) Run() {
	for s := range b.in {
		b.Lock()
		for _, recipient := range b.recipients {
			recipient <- s
		}
		b.Unlock()
	}
	b.stop()
}

// Done returns a channel signalling when broadcasting has stopped
func (b *Broadcaster) Done() <-chan bool {
	return b.done
}

// stop closes broadcast receiver channels and waits for run methods to finish
func (b *Broadcaster) stop() {
	b.Lock()
	defer b.Unlock()
	for _, recipient := range b.recipients {
		close(recipient)
	}
	b.wg.Wait()
	b.done <- true
}

// attach creates and attaches a channel to the broadcaster
func (b *Broadcaster) attach() chan struct{} {
	channel := make(chan struct{})

	b.Lock()
	b.recipients = append(b.recipients, channel)
	b.Unlock()

	return channel
}

// AttachRunner attaches a Run method as broadcast receiver and adds it
// to the waitgroup
func (b *Broadcaster) AttachRunner(runner func(<-chan struct{})) {
	b.wg.Add(1)
	go func() {
		ch := b.attach()
		runner(ch)
		b.wg.Done()
	}()
}
