package main

import (
	"log"
	"micro/pbrpc01/service"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloService struct {
}

func (hs *HelloService) Hello(req *service.Request, resp *service.Response) error {
	resp.Value = "hello: " + req.Value
	return nil
}

var _ service.HelloService = (*HelloService)(nil)

func main() {
	// 注册服务名称
	rpc.RegisterName(service.HelloServiceName, new(HelloService))
	//  建立服务端的tcp链接
	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen err ", err)
		return
	}

	for {
		// 获取每一个链接
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("Accept err", err)
		}
		// 异步rpc处理连接
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
