package main

import (
	"fmt"
	"log"
	"micro/btrpc01/service"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type HelloServiceClient struct {
	*rpc.Client
}

func (hsClient *HelloServiceClient) Hello(req *service.Request, resp *service.Response) error {
	return hsClient.Client.Call(service.HelloServiceName+".Hello", req, resp)
}

var _ service.HelloService = (*HelloServiceClient)(nil)

func DialHelloService(network, address string) (*HelloServiceClient, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Fatal("net.Dial err", err)
		return nil, err
	}
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))
	return &HelloServiceClient{
		client,
	}, nil
}

func main() {
	client, err := DialHelloService("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("Dial err", err)
		return
	}
	resp := &service.Response{}
	err = client.Hello(&service.Request{
		Value: "world",
	}, resp)
	if err != nil {
		return
	}
	fmt.Println(resp)
}
