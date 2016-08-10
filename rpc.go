package mpv

import "net/rpc"

// RPCServer publishes a LLClient over RPC.
type RPCServer struct {
	llclient LLClient
}

// NewRPCServer creates a new RPCServer based on lowlevel client.
func NewRPCServer(client LLClient) *RPCServer {
	return &RPCServer{
		llclient: client,
	}
}

// Exec exposes llclient.Exec via RPC. Not intended to be used directly.
func (s *RPCServer) Exec(args *[]interface{}, res *Response) error {
	resp, err := s.llclient.Exec(*args...)
	*res = *resp
	return err
}

// RPCClient represents a LLClient over RPC.
type RPCClient struct {
	client *rpc.Client
}

// NewRPCClient creates a new RPCClient based on rpc.Client
func NewRPCClient(client *rpc.Client) *RPCClient {
	return &RPCClient{
		client: client,
	}
}

// Exec executes a command over rpc.
func (s *RPCClient) Exec(command ...interface{}) (*Response, error) {
	var res Response
	err := s.client.Call("RPCServer.Exec", &command, &res)
	return &res, err
}
