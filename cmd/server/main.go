package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Kimoto-Norihiro/image-converter/database"
	"github.com/joho/godotenv"

	imageservecepb "github.com/Kimoto-Norihiro/image-converter/pkg/grpc"
	"github.com/Kimoto-Norihiro/image-converter/services/imageservice"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	port := os.Getenv("SERVER_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", 8080))
	if err != nil {
		panic(err)
	}

	db, err := database.NewDB()
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	imageservecepb.RegisterImageServiceServer(s, imageservice.NewImageService(db))
	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
