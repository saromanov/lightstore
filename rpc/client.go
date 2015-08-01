package rpc

import (
	"net"
	"net/rpc"
	"time"
	"log"
)

type RPCClient struct {
	connection *rpc.Client
}

type ClientOptions struct {
	Address string
	Timeout time.Duration
}

func InitClient(opt *ClientOptions) *RPCClient {
	addr := ADDR
	timeout := time.Duration(100) * time.Millisecond

	if opt != nil {
		addr = opt.Address
		timeout = opt.Timeout
	}
	
	connection, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		log.Fatal(err)
	}

	return &RPCClient{connection: rpc.NewClient(connection)}
}

func (client *RPCClient) Get(title string, input,output interface{}) {
	client.connection.Call(title, input, output)
}