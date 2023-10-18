package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	imageservicepb "github.com/Kimoto-Norihiro/image-converter/pkg/grpc"
)

var (
	scanner *bufio.Scanner
	client  imageservicepb.ImageServiceClient
)

func main() {
	fmt.Println("start gRPC Client.")

	scanner = bufio.NewScanner(os.Stdin)

	// 2. gRPCサーバーとのコネクションを確立
	address := "localhost:8080"
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

	// 3. gRPCクライアントを生成
	client = imageservicepb.NewImageServiceClient(conn)

	for {
		fmt.Println("1: image list")
		fmt.Println("2: exit")
		fmt.Print("please enter >")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			ListImages()

		case "2":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}

func ListImages() {
	fmt.Println("image list")
	req := &imageservicepb.ListImagesRequest{}
	res, err := client.ListImages(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
		for i, image := range res.Images {
			fmt.Printf("image %d: %v\n", i, image)
		}
	}
}
