package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"micro/grpc01/service"
)

func main() {
	// grpc.Dial负责和gRPC服务建立连接
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	// 这里会提示，WithInsecure已被弃用，
	// 如果你不想继续使用WithInsecure，可以使用
	// 函数insecure.NewCredentials()返回credentials.TransportCredentials的一个实现。
	// 您可以将其作为DialOption与grpc.WithTransportCredentials一起使用：
	// 但是，API标记为实验性的，因此即使他们已经添加了弃用警告，您也不必立即切换。
	//conn, err := grpc.Dial("localhost:1234",grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Dial err: ", err)
	}
	defer conn.Close()

	// NewHelloServiceClient函数是xxx_grpc.pb.go中自动生成的函数，
	// 基于已经建立的连接构造HelloServiceClient对象,
	// 返回的client其实是一个HelloServiceClient接口对象
	//
	client := service.NewHelloServiceClient(conn)

	// 通过接口定义的方法就可以调用服务端对应gRPC服务提供的方法
	req := &service.Request{Value: "小亮"}
	reply, err := client.Hello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())
}
