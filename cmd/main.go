package main

import (
	"fmt"
	"log"
	"net"

	"github.com/RaniaMidaoui/gomart-authentication-service/pkg/config"
	"github.com/RaniaMidaoui/gomart-authentication-service/pkg/db"
	"github.com/RaniaMidaoui/gomart-authentication-service/pkg/pb"
	"github.com/RaniaMidaoui/gomart-authentication-service/pkg/services"
	"github.com/RaniaMidaoui/gomart-authentication-service/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := db.Init(c.DBUrl)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "gomart-auth-service",
		ExpirationHours: 24 * 365,
	}

	lis, err := net.Listen("tcp", c.Port)

	if err != nil {
		log.Fatalln("Failed to listing:", err)
	}

	fmt.Println("Mini Project DevOps!!")
	fmt.Println("Auth Svc on", c.Port)

	s := services.Server{
		H:   h,
		Jwt: jwt,
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalln("Failed to serve:", err)
	}
}
