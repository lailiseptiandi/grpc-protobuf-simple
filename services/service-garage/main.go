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

var localstorage *model.GarageListByUser

func init() {
	localstorage = new(model.GarageListByUser)
	localstorage.List = make(map[string]*model.GarageList)
}

type GaragesServer struct{}

func (GaragesServer) Add(ctx context.Context, param *model.GarageAndUserId) (*empty.Empty, error) {
	userId := param.UserId
	garage := param.Garage

	if _, ok := localstorage.List[userId]; !ok {
		localstorage.List[userId] = new(model.GarageList)
		localstorage.List[userId].List = make([]*model.Garage, 0)
	}

	localstorage.List[userId].List = append(localstorage.List[userId].List, garage)
	log.Println("Adding garage", garage.String(), "for user", userId)

	return new(empty.Empty), nil
}

func (GaragesServer) List(ctx context.Context, param *model.GarageUserId) (*model.GarageList, error) {
	userId := param.UserId

	return localstorage.List[userId], nil
}

func main() {
	srv := grpc.NewServer()
	var garageSrv GaragesServer
	model.RegisterGaragesServer(srv, garageSrv)

	log.Println("Starting RPC server at", config.SERVICE_GARAGE_PORT)

	l, err := net.Listen("tcp", config.SERVICE_GARAGE_PORT)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", config.SERVICE_GARAGE_PORT, err)
	}

	log.Fatal(srv.Serve(l))
}
