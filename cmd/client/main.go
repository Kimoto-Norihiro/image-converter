package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	imageservicepb "github.com/Kimoto-Norihiro/image-converter/pkg/grpc"
	"github.com/joho/godotenv"
)

var (
	scanner *bufio.Scanner
	client  imageservicepb.ImageServiceClient
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	fmt.Println("start gRPC Client.")
	scanner = bufio.NewScanner(os.Stdin)

	port := os.Getenv("SERVER_PORT")
	address := fmt.Sprintf("localhost:%v", port)
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatal("Connection failed.")
		return
	}
	defer conn.Close()

	client = imageservicepb.NewImageServiceClient(conn)

	for {
		fmt.Println("1: create image")
		fmt.Println("2: image list")
		fmt.Println("3: convert images")
		fmt.Println("4: exit")
		fmt.Print("please enter > ")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			CreateImage()
		case "2":
			ConvertImages()
		case "3":
			ListImages()
		case "4":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}
