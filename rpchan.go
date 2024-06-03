// Package provides channel-like semantics over a TCP connection using net/rpc.
package rpchan

import (
	"errors"
	"net"
	"net/rpc"

	"github.com/lucafmarques/rpchan/internal"
)

// RPChan
type RPChan[T any] struct {
	setupC   func() error
	setupR   func() error
	client   *rpc.Client
	listener *net.Listener
	receiver *internal.Receiver[T]
}

// Send implements Go channels' send operation.
//
// A call to Send, much like sending over a Go channel, may block.
// The first call to Send may error on dialing the TCP address of the RPChan.
//
// Since this involves a network call, Send can return an error.
func (ch *RPChan[T]) Send(v T) error {
	if err := ch.setupC(); err != nil {
		return err
	}

	return ch.client.Call("Channel.Send", v, nil)
}

// Receive implements Go channels' receive operation.
//
// A call to Receive, much like receiving over a Go channel, may block.
// The first call to Receive may error on listening on the TCP address of RPChan.
//
// Since this may involve a network call, Receive can return an error.
func (ch *RPChan[T]) Receive() (v T, err error) {
	if err = ch.setupR(); err != nil {
		return v, err
	}

	r, ok := <-ch.receiver.Channel
	if !ok {
		return v, net.ErrClosed
	}

	return r, err
}

// Listen implements a GOEXPERIMENT=rangefunc iterator.
//
// When used in a for-range loop it works exacly like a Go channel.
func (ch *RPChan[T]) Listen() func(func(T, error) bool) {
	return func(yield func(T, error) bool) {
		for {
			if v, err := ch.Receive(); !yield(v, err) || errors.Is(err, net.ErrClosed) {
				return
			}
		}
	}
}

// Close implements the close built-in.
//
// Since closing involves I/O, it can return an error containing
// the RPC client's Close() error and/or the TCP listener Close() error.
func (ch *RPChan[T]) Close() error {
	var errs []error

	if ch.client != nil {
		errs = append(errs, ch.client.Call("Channel.Close", 0, nil), ch.client.Close())
	}
	if ch.receiver != nil {
		errs = append(errs, (*(ch.listener)).Close(), ch.receiver.Close(0, nil))
	}

	return errors.Join(errs...)
}

// New creates an RPChan[T] over addr, with an optional buffer, and returns a reference to it.
//
// The returned RPChan[T] will not start a client nor a server unless their
// related methods are called, [RPChan.Send] and [RPChan.Receive] or [RPChan.Listen], respectively.
func New[T any](addr string, buffer ...uint) *RPChan[T] {
	var bufsize uint
	if len(buffer) > 0 {
		bufsize = buffer[0]
	}

	ch := &RPChan[T]{}
	ch.setupC = func() error {
		if ch.client == nil {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				return err
			}

			ch.client = rpc.NewClient(conn)
		}

		return nil
	}
	ch.setupR = func() error {
		if ch.receiver == nil {
			srv := rpc.NewServer()
			rec := internal.NewReceiver[T](bufsize)
			srv.RegisterName("Channel", rec)

			list, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}

			go srv.Accept(list)

			ch.receiver, ch.listener = rec, &list
		}

		return nil
	}

	return ch
}
