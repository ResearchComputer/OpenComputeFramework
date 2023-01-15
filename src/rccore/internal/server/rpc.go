package server

import (
	rpc "github.com/ethereum/go-ethereum/rpc"
)

func NewRPCServer() *rpc.Server {
	workerService := new(WorkerService)
	rpcServer := rpc.NewServer()
	err := rpcServer.RegisterName("worker", workerService)
	if err != nil {
		panic(err)
	}
	return rpcServer
}
