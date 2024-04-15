package rpchan_test

import (
	"net"
	"testing"
	"time"

	"github.com/carlmjohnson/be"
	"github.com/lucafmarques/rpchan"
)

var TestTable = map[string]func(t *testing.T){
	"SEND_AND_CLOSE": func(t *testing.T) {
		sender := rpchan.New[int](":9760")
		receiver := rpchan.New[int](":9760")

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := sender.Send(10)
			be.NilErr(t, err)

			err = sender.Close()
			be.NilErr(t, err)
		}()

		val, err := receiver.Receive()
		be.NilErr(t, err)
		be.Equal(t, 10, val)
	},
	"SEND_AND_LISTEN": func(t *testing.T) {
		sender := rpchan.New[int](":9761")
		receiver := rpchan.New[int](":9761")
		vals := []int{1, 2, 3}

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := sender.Send(vals[0])
			be.NilErr(t, err)

			err = sender.Send(vals[1])
			be.NilErr(t, err)

			err = sender.Send(vals[2])
			be.NilErr(t, err)

			err = sender.Close()
			be.NilErr(t, err)
		}()

		i := 0
		for val, err := range receiver.Listen() {
			if i < len(vals) {
				be.Equal(t, vals[i], val)
				be.NilErr(t, err)
			}
			i++
		}
	},
	"RECEIVE_AND_CLOSE": func(t *testing.T) {
		sender := rpchan.New[int](":9762")
		receiver := rpchan.New[int](":9762")

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := sender.Send(10)
			be.NilErr(t, err)
		}()

		val, err := receiver.Receive()
		be.NilErr(t, err)
		be.Equal(t, 10, val)

		err = receiver.Close()
		be.NilErr(t, err)
	},
	"RECEIVE_AFTER_CLOSE": func(t *testing.T) {
		sender := rpchan.New[int](":9762")
		receiver := rpchan.New[int](":9762")

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := sender.Send(10)
			be.NilErr(t, err)
		}()

		val, err := receiver.Receive()
		be.NilErr(t, err)
		be.Equal(t, 10, val)

		err = receiver.Close()
		be.NilErr(t, err)

		_, err = receiver.Receive()
		be.Equal(t, net.ErrClosed.Error(), err.Error())
	},
	"BUFFERED_RECEIVE_AND_CLOSE": func(t *testing.T) {
		sender := rpchan.New[int](":9763")
		receiver := rpchan.New[int](":9763", 100)

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := sender.Send(10)
			be.NilErr(t, err)
		}()

		val, err := receiver.Receive()
		be.NilErr(t, err)
		be.Equal(t, 10, val)

		err = receiver.Close()
		be.NilErr(t, err)
	},
	"SEND_ADDRESS_IN_USE": func(t *testing.T) {
		receiver1 := rpchan.New[int](":9764")
		receiver2 := rpchan.New[int](":9764")

		go func() {
			time.Sleep(10 * time.Millisecond)
			_, err := receiver1.Receive()
			be.Equal(t, "listen tcp :9764: bind: address already in use", err.Error())
		}()

		err := receiver1.Send(10)
		if err != nil {
			be.Equal(t, "dial tcp :9764: connect: connection refused", err.Error())

		}

		err = receiver2.Send(10)
		if err != nil {
			be.Equal(t, "dial tcp :9764: connect: connection refused", err.Error())
		}

		err = receiver1.Close()
		be.NilErr(t, err)
	},
	"RECEIVE_ADDRESS_IN_USE": func(t *testing.T) {
		receiver1 := rpchan.New[int](":9765")
		receiver2 := rpchan.New[int](":9765")

		go func() {
			time.Sleep(10 * time.Millisecond)
			_, err := receiver1.Receive()
			be.Equal(t, "listen tcp :9765: bind: address already in use", err.Error())
		}()

		go func() {
			time.Sleep(10 * time.Millisecond)
			err := receiver1.Send(10)
			be.NilErr(t, err)
		}()

		_, err := receiver2.Receive()
		be.NilErr(t, err)

		err = receiver1.Close()
		be.NilErr(t, err)
	},
}

func TestRPChan(t *testing.T) {
	for test, f := range TestTable {
		t.Run(test, f)
	}
}
