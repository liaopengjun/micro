package main

import (
	"fmt"
	"log"
	"micro/rpc/service"
	"net/rpc"
)

type HelloServiceClient struct {
	*rpc.Client
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(service.HelloServiceName+".Hello", request, reply)
}

// 静态检查，同上面一样
var _ service.HelloService = (*HelloServiceClient)(nil)

// 通过rpc.Dial拨号RPC服务，建立连接,并将获取连接后的客户端返回
func DialHelloService(network, address string) (*HelloServiceClient, error) {
	client, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &HelloServiceClient{client}, nil
}

func main() {
	client, err := DialHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatalln("dialing:", err)
	}
	var reply string
	// 在使用goland的时候就会提示
	err = client.Hello("world", &reply)
	fmt.Println(reply)
}
