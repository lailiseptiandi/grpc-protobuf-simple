package main

import (
	"context"
	"grpc-protobuf/common/config"
	"grpc-protobuf/common/model"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var localstorage *model.UserList

func init() {
	localstorage = new(model.UserList)
	localstorage.List = make([]*model.User, 0)
}

type UsersServer struct{}

func (UsersServer) Register(ctx context.Context, param *model.User) (*empty.Empty, error) {
	localstorage.List = append(localstorage.List, param)

	log.Println("Registering user", param.String())

	return new(empty.Empty), nil
}

func (UsersServer) List(ctx context.Context, void *empty.Empty) (*model.UserList, error) {
	return localstorage, nil
}

func main() {
	srv := grpc.NewServer()
	var userSrv UsersServer
	model.RegisterUsersServer(srv, userSrv)

	log.Println("Starting RPC server at", config.SERVICE_USER_PORT)

	l, err := net.Listen("tcp", config.SERVICE_USER_PORT)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", config.SERVICE_USER_PORT, err)
	}

	log.Fatal(srv.Serve(l))
}
