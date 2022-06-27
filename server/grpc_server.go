package main

import (
	"fmt"
	"item-admin/proto"
	"item-admin/service"
	"net"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listen, err := net.Listen("tcp", ":8001")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()
	proto.RegisterItemServiceServer(s, &service.Server{})
	reflection.Register(s)

	defer func() {
		s.Stop()
		listen.Close()
	}()

	fmt.Println("Serving 8001...")
	err = s.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
