# `rpchan`: channel-like semantics over net/rpc
![coverage](https://img.shields.io/badge/coverage-100.0%25-brightgreen)
[![Go Reference](https://pkg.go.dev/badge/github.com/lucafmarques/rpchan.svg)](https://pkg.go.dev/github.com/lucafmarques/rpchan)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucafmarques/rpchan)](https://goreportcard.com/report/github.com/lucafmarques/rpchan)

`rpchan` provides channel-like semantics over a TCP connection using `net/rpc`.

```
go get github.com/lucafmarques/rpchan
```

It achieves this by providing a minimal API on the `RPChan[T any]` type: 
- Use [`Send`](https://pkg.go.dev/github.com/lucafmarques/rpchan#RPChan.Send) if you want to send a `T`, similarly to `ch <- T`
- Use [`Receive`](https://pkg.go.dev/github.com/lucafmarques/rpchan#RPChan.Receive) if you want to receive a `T`, similarly to `<-ch`
- Use [`Close`](https://pkg.go.dev/github.com/lucafmarques/rpchan#RPChan.Close) if you want to close the channel, similarly to `close(ch)`
- Use [`Listen`](https://pkg.go.dev/github.com/lucafmarques/rpchan#RPChan.Listen) if want to [iterate](#rangefunc) on the channel, similarly to `for v := range ch`

Those four methods are enough to mimic one-way send/receive channel-like semantics. 

Be mindful that since network calls are involved, error returns are needed to allow callers to react to network errors.

It's advisable, but not mandatory, to use the same type on both the receiver and sender. This is because `rpchan` follows the [`encoding/gob`](https://pkg.go.dev/encoding/gob#hdr-Types_and_Values) guidelines for encoding types and values.

## Examples
<table>
<tr>
<th><code>sender.go</code></th>
<th><code>receiver.go</code></th>
</tr>
<tr>
<td>
  
```go
package main

import (
    "github.com/lucafmarques/rpchan"
)

func main() {
    ch := rpchan.New[int](":9091")
    err := ch.Send(20)
    // ... error handling
    err = ch.Close()
}
```
</td>
<td>
  
```go
package main

import (
    "github.com/lucafmarques/rpchan"
)

func main() {
    ch := rpchan.New[int](":9091", 100)
    for v, err := range ch.Listen() {
        // ... error handling + use v
    }
}
```
</td>
</tr>
</table>



<details><summary>## Tasks</summary>

Executable `make`-like tasks using [`xc`](https://github.com/joerdav/xc).

### test
```
GOEXPERIMENT=rangefunc go test ./...
```

### tag

Creates new major|minor|patch tag.

requires: test  
inputs: VERSION  
env: VERSION=patch  

```
git fetch --tags || true
CURR_VERSION=`git describe --abbrev=0 --tags 2>/dev/null`
SPLITS=(${CURR_VERSION//./ })
MAJOR=${SPLITS[0]//v}
MINOR=${SPLITS[1]}
PATCH=${SPLITS[2]}

case $VERSION in
    major)
        MAJOR=$((MAJOR+1))
        MINOR=0
        PATCH=0
        ;;
    minor)
        MINOR=$((MINOR+1))
        PATCH=0
        ;;
    patch)
        PATCH=$((PATCH+1))
        ;;
esac 

TAG="v$MAJOR.$MINOR.$PATCH"
git tag $TAG -m "tag($VERSION): Release version $TAG"
```
</details>

---

## rangefunc
If built with Go 1.22 and `GOEXPERIMENT=rangefunc`, the `Listen` method can be used on a for-range loop, working exactly like a Go channel would.

---

README.md inspired by [`sourcegraph/conc`](https://github.com/sourcegraph/conc)
