package internal

import "net"

// Receiver needs to be exported to be used as an RPC handler.
type Receiver[T any] struct {
	Channel chan T
}

func NewReceiver[T any](buf uint) *Receiver[T] { return &Receiver[T]{make(chan T, buf)} }

// Send implements the function signature for an RPC handler.
func (r *Receiver[T]) Send(item T, _ *bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// safe to assume that we're recovering from a send on a closed channel
			err = net.ErrClosed
		}
	}()

	r.Channel <- item
	return nil
}

// Close implements the function signature for an RPC handler.
// It handles a client closing the communication over the rpchan.
func (r *Receiver[T]) Close(_ int, _ *bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			// safe to assume that we're recovering from a close on a closed channel
			err = net.ErrClosed
		}
	}()
	close(r.Channel)
	return nil
}
