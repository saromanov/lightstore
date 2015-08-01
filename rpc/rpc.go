package rpc

import (
    "net"
	"net/rpc"
	"log"
)

type RPCData struct {
	address string
}

const (
	ADDR = ":9873"
)

func RegisterRPCFunction(obj interface{}){
	rpc.Register(obj)
}

func Init(address string) *RPCData {
	rpcdata := new(RPCData)
	if address == "" {
		address = ADDR
	}

	rpcdata.address = address
	return rpcdata
}

func (rpcdata *RPCData) Run() {
	go func() {
		l, e := net.Listen("tcp", rpcdata.address)
		if e != nil {
			log.Fatal("listen error:", e)
		}

		rpc.Accept(l)
	}()
}
