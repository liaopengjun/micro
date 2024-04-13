package service

const HelloServiceName = "HelloService"

type HelloService interface {
	Hello(request string, reply *string) error
}
