package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"micro/grpc02/service"
)

func unAry(client service.EchoClient) {
	resp, err := client.UnaryEcho(context.Background(), &service.EchoRequest{Message: "hello world21"})
	if err != nil {
		fmt.Printf("client unary echo err:%v\n", err)
	}
	fmt.Printf("client unary echo resp:%v\n", resp)
}

/*
	serveStream 接收服务端返回多个流数据

1. 建立连接 获取client
2. 通过 client 获取stream
3. for循环中通过stream.Recv()依次获取服务端推送的消息
4. err==io.EOF则表示服务端关闭stream了
*/
func serveStream(client service.EchoClient) {
	// 调用获取stream
	stream, err := client.ServerStreamingEcho(context.Background(), &service.EchoRequest{
		Message: "hello world",
	})
	// for循环获取服务端推送的消息
	if err != nil {
		log.Fatalf("server streaming echo err:%v\n", err)
		return
	}

	for {
		// 获取服务端send过来的数据
		resp, err := stream.Recv()
		if err == io.EOF {
			// 读取完了
			break
		}
		if err != nil {
			log.Fatalf("recv echo err:%v\n", err)
			continue
		}
		// 打印服务端send过来的数据
		fmt.Printf("server receive resp:%v\n", resp.GetMessage())
	}
}

func main() {
	// 建立grpc链接
	conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("client dial err:%v\n", err)
		return
	}
	defer conn.Close()
	// 调用pb生成client
	client := service.NewEchoClient(conn)
	// 调用grpc提供函数方法
	//unAry(client)
	serveStream(client)
}
