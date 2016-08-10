package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/blang/mpv"
)

func main() {
	ll := mpv.NewIPCClient("/tmp/mpvsocket")
	s := mpv.NewRPCServer(ll)
	rpc.Register(s)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":9999")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}
	go http.Serve(l, nil)

	// Client
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:9999")
	if err != nil {
		log.Fatal("Listen error: ", err)
	}

	rpcc := mpv.NewRPCClient(client)
	c := mpv.NewClient(rpcc)
	c.SetFullscreen(false)
	c.SetPause(true)

}
