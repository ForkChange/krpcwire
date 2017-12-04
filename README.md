# krpcwire

Low level implementation of KRPC network layer of DHT.

## Install
```bash
go get "github.com/fanpei91/krpcwire"
```

## Usage
```go
import "github.com/fanpei91/krpcwire"
```

#### func  OnInQuery

```go
func OnInQuery(on func(query InRequest, from net.UDPAddr)) option
```

#### func  Timeout

```go
func Timeout(t time.Duration) option
```

#### func  TransIDSize

```go
func TransIDSize(n int) option
```

#### type InRequest

```go
type InRequest map[string]interface{}
```


#### type InResponse

```go
type InResponse map[string]interface{}
```


#### type OutError

```go
type OutError []interface{}
```


#### type OutQuery

```go
type OutQuery struct {
	Q string
	A map[string]interface{}
}
```


#### type OutQueryCallback

```go
type OutQueryCallback func(req *OutRequest, res InResponse, timeout bool, from net.UDPAddr)
```


#### type OutRequest

```go
type OutRequest struct {
	OutQuery
	Tid TransID
	Y   string
}
```


#### type OutResponse

```go
type OutResponse map[string]interface{}
```


#### type TransID

```go
type TransID string
```


#### type Wire

```go
type Wire struct {
}
```


#### func  NewWire

```go
func NewWire(socket *net.UDPConn, options ...option) *Wire
```

#### func (*Wire) Cancel

```go
func (w *Wire) Cancel(tid TransID)
```

#### func (*Wire) Error

```go
func (w *Wire) Error(req *OutRequest, err OutError, to net.UDPAddr)
```

#### func (*Wire) Query

```go
func (w *Wire) Query(query OutQuery, cb OutQueryCallback, to net.UDPAddr) TransID
```

#### func (*Wire) Reply

```go
func (w *Wire) Reply(req *OutRequest, res OutResponse, to net.UDPAddr)
```
