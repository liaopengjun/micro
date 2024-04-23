package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"micro/grpc02/service"
	"strconv"
	"sync"
	"time"
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

// clientStream 客户端流
/**
1.建立连接并获取client
2.获取stream并通过send不断想服务端发送数据
3.发送完成后通过stream.CloseAndRecv() 关闭steam并接收服务端返回结果
*/
func clientStream(client service.EchoClient) {
	stream, err := client.ClientStreamingEcho(context.Background())
	if err != nil {
		log.Fatalf("client streaming echo err:%v\n", err)
	}
	// 发送数据
	for i := 0; i < 100; i++ {
		err := stream.Send(&service.EchoRequest{Message: "hello " + strconv.Itoa(i)})
		if err != nil {
			log.Fatalf("send echo err:%v\n", err)
			continue
		}
	}
	// 服务端返回结果
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CloseAndRecv err:%v\n", err)
	}
	log.Printf("resp:%v\n", resp.GetMessage())

}

// bidirectionalStream 双向流
/*
1. 建立连接 获取client
2. 通过client获取stream
3. 开两个goroutine 分别用于Recv()和Send()
	3.1 一直Recv()到err==io.EOF(即服务端关闭stream)
	3.2 Send()则由自己控制
4. 发送完毕调用 stream.CloseSend()关闭stream 必须调用关闭 否则Server会一直尝试接收数据 一直报错...
*/
func bidirectionalStream(client service.EchoClient) {
	var wg sync.WaitGroup
	// 2. 调用方法获取stream
	stream, err := client.BidirectionalStreamingEcho(context.Background())
	if err != nil {
		log.Fatalf("client streaming echo err:%v\n", err)
		panic(err)
	}
	// 3.开两个goroutine 分别用于Recv()和Send()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("Server Closed")
				break
			}
			if err != nil {
				continue
			}
			fmt.Printf("Recv Data:%v \n", req.GetMessage())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 2; i++ {
			err := stream.Send(&service.EchoRequest{Message: "hello world"})
			if err != nil {
				log.Printf("send error:%v\n", err)
			}
			time.Sleep(time.Second)
		}
		// 4. 发送完毕关闭stream
		err := stream.CloseSend()
		if err != nil {
			log.Printf("Send error:%v\n", err)
			return
		}
	}()
	wg.Wait()
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
	//serveStream(client)
	//clientStream(client)
	bidirectionalStream(client)
}
