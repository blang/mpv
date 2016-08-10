package mpv

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"
)

// Response received from mpv. Can be an event or a user requested response.
type Response struct {
	Err       string      `json:"error"`
	Data      interface{} `json:"data"` // May contain float64, bool or string
	Event     string      `json:"event"`
	RequestID int         `json:"request_id"`
}

// request sent to mpv. Includes request_id for mapping the response.
type request struct {
	Command   []interface{}  `json:"command"`
	RequestID int            `json:"request_id"`
	Response  chan *Response `json:"-"`
}

func newRequest(cmd ...interface{}) *request {
	return &request{
		Command:   cmd,
		RequestID: rand.Intn(10000),
		Response:  make(chan *Response, 1),
	}
}

// LLClient is the most low level interface
type LLClient interface {
	Exec(command ...interface{}) (*Response, error)
}

// IPCClient is a low-level IPC client to communicate with the mpv player via socket.
type IPCClient struct {
	socket  string
	timeout time.Duration
	comm    chan *request

	mu     sync.Mutex
	reqMap map[int]*request // Maps RequestIDs to Requests for response association
}

// NewIPCClient creates a new IPCClient connected to the given socket.
func NewIPCClient(socket string) *IPCClient {
	c := &IPCClient{
		socket:  socket,
		timeout: 2 * time.Second,
		comm:    make(chan *request),
		reqMap:  make(map[int]*request),
	}
	c.run()
	return c
}

// dispatch dispatches responses to the corresponding request
func (c *IPCClient) dispatch(resp *Response) {
	if resp.Event == "" { // No Event
		c.mu.Lock()
		defer c.mu.Unlock()
		if req, ok := c.reqMap[resp.RequestID]; ok { // Lookup requestID in request map
			delete(c.reqMap, resp.RequestID)
			req.Response <- resp
			return
		}
		// Discard response
	} else { // Event
		// TODO: Implement Event support
	}
}

func (c *IPCClient) run() {
	conn, err := net.Dial("unix", c.socket)
	if err != nil {
		panic(err)
	}
	go c.readloop(conn)
	go c.writeloop(conn)
	// TODO: Close connection
}

func (c *IPCClient) writeloop(conn io.Writer) {
	for {
		req, ok := <-c.comm
		if !ok {
			panic("Communication channel closed")
		}
		b, err := json.Marshal(req)
		if err != nil {
			// TODO: Discard request, maybe send error downstream
			// log.Printf("Discard request %v with error: %s", req, err)
			continue
		}
		c.mu.Lock()
		c.reqMap[req.RequestID] = req
		c.mu.Unlock()
		b = append(b, '\n')
		_, err = conn.Write(b)
		if err != nil {
			// TODO: Discard request, maybe send error downstream
			// TODO: Remove from reqMap?
		}
	}
}

func (c *IPCClient) readloop(conn io.Reader) {
	rd := bufio.NewReader(conn)
	for {
		data, err := rd.ReadBytes('\n')
		if err != nil {
			// TODO: Handle error
			continue
		}
		var resp Response
		err = json.Unmarshal(data, &resp)
		if err != nil {
			// TODO: Handle error
			continue
		}
		c.dispatch(&resp)
	}
}

// Timeout errors while communicating via IPC
var (
	ErrTimeoutSend = errors.New("Timeout while sending command")
	ErrTimeoutRecv = errors.New("Timeout while receiving response")
)

// Exec executes a command via ipc and returns the response.
// A request can timeout while sending or while waiting for the response.
// An error is only returned if there was an error in the communication.
// The client has to check for `response.Error` in case the server returned
// an error.
func (c *IPCClient) Exec(command ...interface{}) (*Response, error) {
	req := newRequest(command...)
	select {
	case c.comm <- req:
	case <-time.After(c.timeout):
		return nil, ErrTimeoutSend
	}

	select {
	case res, ok := <-req.Response:
		if !ok {
			panic("Response channel closed")
		}
		return res, nil
	case <-time.After(c.timeout):
		return nil, ErrTimeoutRecv
	}
}
