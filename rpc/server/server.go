package server

import (
	"log"
	"micro/rpc/service"
	"net"
	"net/rpc"
)

type HelloServer struct{}

func (p *HelloServer) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

// 通过接口约束HelloService服务
var _ service.HelloService = (*HelloServer)(nil)

func Start() {
	rpc.RegisterName(service.HelloServiceName, new(HelloServer))
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}
	log.Println("启动rpc")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}
		go rpc.ServeConn(conn)
	}
}
