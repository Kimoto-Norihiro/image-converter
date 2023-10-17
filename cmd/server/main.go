package main

import (
	// (一部抜粋)
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	dns := "n000r111:password@tcp(mysql:3306)/image_converter?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := database.NewDB(dns)
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
