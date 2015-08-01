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

func (client *RPCClient) Init(opt *ClientOptions) *RPCClient {
	addr := opt.Address
	if addr == "" {
		addr = ADDR
	}

	timeout := opt.Timeout
	if timeout == 0 {
		timeout = 2
	}

	connection, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		log.Fatal(err)
	}

	return &RPCClient{connection: rpc.NewClient(connection)}
}

func (client *RPCClient) Get(procname string) {
	//client.connection.Call()
}
