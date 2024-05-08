package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"micro/grpc02/service"
	"time"
)

func main() {

	retryPolicy := `{
		"methodConfig": [{
		  "name": [{"service": "service.Echo","method":"UnaryEcho"}],
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`

	conn, err := grpc.Dial("127.0.0.1:1234", grpc.WithInsecure(), grpc.WithDefaultServiceConfig(retryPolicy))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()
	c := service.NewEchoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	reply, err := c.UnaryEcho(ctx, &service.EchoRequest{Message: "Try and Success"})
	if err != nil {
		log.Fatalf("UnaryEcho error: %v", err)
	}
	log.Printf("UnaryEcho reply: %v", reply)

}
