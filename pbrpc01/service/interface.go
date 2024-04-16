package service

const HelloServiceName = "HelloService"

type HelloService interface {
	// Hello
	// 这里的 Request 和  Response 是基于protobuf生成的service.pb.go里的结构
	Hello(request *Request, response *Response) error
}
