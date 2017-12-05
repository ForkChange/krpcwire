/*
 * Low level implementation of KRPC network layer of DHT.
 */

package krpcwire

import (
	"bytes"
	"crypto/rand"
	"net"
	"sync"
	"time"

	"github.com/fanpei91/bencode"
)

type OutQueryCallback func(req *OutRequest, res InResponse, timeout bool, from net.UDPAddr)
type OutResponse map[string]interface{}
type OutError []interface{}
type OutQuery struct {
	Q string
	A map[string]interface{}
}
type OutRequest struct {
	OutQuery
	TransID string
	Y       string
}

type InRequest map[string]interface{}
type InResponse map[string]interface{}

type option func(*Wire)

type request struct {
	createdAt time.Time
	to        net.UDPAddr
	callback  OutQueryCallback
	message   *OutRequest
}

func randomTransID(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return string(b)
}

func OnInQuery(on func(query InRequest, from net.UDPAddr)) option {
	return func(w *Wire) {
		w.onInQuery = on
	}
}

func Timeout(t time.Duration) option {
	return func(w *Wire) {
		w.timeout = t
	}
}

func TransIDSize(n int) option {
	return func(w *Wire) {
		w.transIDSize = n
	}
}

type Wire struct {
	socket      *net.UDPConn
	timeout     time.Duration
	reqs        sync.Map
	transIDSize int
	onInQuery   func(query InRequest, from net.UDPAddr)
}

func NewWire(socket *net.UDPConn, options ...option) *Wire {
	w := &Wire{
		socket:      socket,
		timeout:     2 * time.Second,
		transIDSize: 3,
	}
	for _, option := range options {
		option(w)
	}
	go w.listen()
	go w.sweep()
	return w
}

func (w *Wire) Query(query OutQuery, cb OutQueryCallback, to net.UDPAddr) (transID string) {
	transID = randomTransID(w.transIDSize)
	message := &OutRequest{
		OutQuery: query,
		TransID:  transID,
		Y:        "q",
	}
	req := &request{
		createdAt: time.Now(),
		to:        to,
		callback:  cb,
		message:   message,
	}
	w.send(map[string]interface{}{
		"t": transID,
		"y": message.Y,
		"q": message.Q,
		"a": message.A,
	}, to)
	w.reqs.Store(transID, req)
	return transID
}

func (w *Wire) Reply(transID string, res OutResponse, to net.UDPAddr) {
	w.send(map[string]interface{}{
		"t": transID,
		"y": "r",
		"r": map[string]interface{}(res),
	}, to)
}

func (w *Wire) Error(transID string, err OutError, to net.UDPAddr) {
	w.send(map[string]interface{}{
		"t": transID,
		"y": "e",
		"e": []interface{}(err),
	}, to)
}

func (w *Wire) Cancel(transID string) {
	w.reqs.Delete(transID)
}

func (w *Wire) listen() {
	buf := make([]byte, 8192)
	for {
		if n, from, err := w.socket.ReadFromUDP(buf); err == nil {
			go w.onMessage(buf[:n], *from)
		}
	}
}

func (w *Wire) sweep() {
	for range time.Tick(w.timeout) {
		now := time.Now()
		w.reqs.Range(func(key, value interface{}) bool {
			if req, ok := value.(*request); ok {
				if now.Sub(req.createdAt) > w.timeout {
					w.reqs.Delete(key)
					go req.callback(req.message, nil, true, req.to)
				}
			}
			return true
		})
	}
}

func (w *Wire) onMessage(data []byte, from net.UDPAddr) {
	dict, err := bencode.Decode(bytes.NewBuffer(data))
	if err != nil {
		return
	}
	y, ok := dict["y"].(string)
	if !ok {
		return
	}
	switch y {
	case "q":
		if w.onInQuery != nil {
			w.onInQuery(dict, from)
		}
	case "r", "e":
		w.onResponse(dict, from)
	}
}

func (w *Wire) onResponse(res InResponse, from net.UDPAddr) {
	t, ok := res["t"].(string)
	if !ok {
		return
	}
	v, ok := w.reqs.Load(t)
	if !ok {
		return
	}
	w.reqs.Delete(t)
	req := v.(*request)
	now := time.Now()
	if now.Sub(req.createdAt) > w.timeout {
		req.callback(req.message, nil, true, from)
		return
	}
	req.callback(req.message, res, false, from)
}

func (w *Wire) send(msg map[string]interface{}, to net.UDPAddr) {
	w.socket.WriteToUDP(bencode.Encode(msg), &to)
}
