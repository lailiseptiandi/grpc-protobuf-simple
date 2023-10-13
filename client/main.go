package main

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc-protobuf/common/config"
	"grpc-protobuf/common/model"
	"log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

func serviceGarage() model.GaragesClient {
	port := config.SERVICE_GARAGE_PORT
	conn, err := grpc.Dial(port, grpc.WithInsecure())

	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return model.NewGaragesClient(conn)
}

func serviceUser() model.UsersClient {
	port := config.SERVICE_USER_PORT
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return model.NewUsersClient(conn)
}

func main() {
	userTest()
	garageTest()
}

func userTest() {
	user1 := model.User{
		Id:       "1",
		Name:     "laili septiandi",
		Password: "hashpassword",
		Gender:   model.UserGender_MALE,
	}

	user2 := model.User{
		Id:       "2",
		Name:     "kang jamil",
		Password: "hashpassword",
		Gender:   model.UserGender_MALE,
	}

	user := serviceUser()

	fmt.Println("===> user test  <===")
	user.Register(context.Background(), &user1)
	user.Register(context.Background(), &user2)

	res1, err := user.List(context.Background(), new(empty.Empty))
	if err != nil {
		log.Fatal(err.Error())
	}
	res1String, _ := json.Marshal(res1.List)
	log.Println(string(res1String))
}

func garageTest() {
	garage := serviceGarage()
	user := serviceUser()
	user1, err := user.List(context.Background(), new(empty.Empty))
	if err != nil {
		log.Fatal(err.Error())
	}
	payload := model.GarageAndUserId{
		UserId: user1.List[0].Id,
		Garage: &model.Garage{
			Id:   "1",
			Name: "Indonesia",
			Coordinate: &model.GarageCoordinate{
				Latitude:  95.9614608,
				Longitude: -2.2243723,
			},
		},
	}
	garage.Add(context.Background(), &payload)

	res1, err := garage.List(context.Background(), &model.GarageUserId{
		UserId: payload.UserId,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	res1String, _ := json.Marshal(res1.List)

	log.Println("Hasil List Garage :", string(res1String))
}
