package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"micro/grpc02/service"
	"time"
)

func unaryCall(c service.EchoClient, requestID int, message string, want codes.Code) {
	// 每次都指定1秒超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &service.EchoRequest{Message: message}
	msg, err := c.UnaryEcho(ctx, req)
	got := status.Code(err)
	fmt.Printf("[%v] wanted = %v, got = %v, msg = %v \n ", requestID, want, got, msg)
}

func main() {
	c, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	defer c.Close()
	client := service.NewEchoClient(c)
	unaryCall(client, 1, "Normal", codes.InvalidArgument)
}
