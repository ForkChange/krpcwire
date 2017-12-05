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


#### type OnReply

```go
type OnReply func(req *OutRequest, res InResponse, timeout bool, from net.UDPAddr)
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


#### type OutRequest

```go
type OutRequest struct {
	OutQuery
	TransID string
	Y       string
}
```


#### type OutResponse

```go
type OutResponse map[string]interface{}
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
func (w *Wire) Cancel(transID string)
```

#### func (*Wire) Error

```go
func (w *Wire) Error(transID string, err OutError, to net.UDPAddr)
```

#### func (*Wire) Query

```go
func (w *Wire) Query(query OutQuery, cb OnReply, to net.UDPAddr) (transID string)
```

#### func (*Wire) Reply

```go
func (w *Wire) Reply(transID string, res OutResponse, to net.UDPAddr)
```
