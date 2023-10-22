package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

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
		fmt.Println("1: image list")
		fmt.Println("2: convert images")
		fmt.Println("3: create image")
		fmt.Println("4: exit")
		fmt.Print("please enter > ")

		scanner.Scan()
		in := scanner.Text()

		switch in {
		case "1":
			ListImages()

		case "2":
			ConvertImages()

		case "3":
			CreateImage()

		case "4":
			fmt.Println("bye.")
			goto M
		}
	}
M:
}

func ListImages() {
	fmt.Println("show images")
	req := &imageservicepb.ListImagesRequest{}
	res, err := client.ListImages(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
		for i, image := range res.Images {
			fmt.Printf("image %d: %v\n", i+1, image)
		}
	}
}

func ConvertImages() {
	fmt.Println("convert images")
	req := &imageservicepb.ConvertImagesRequest{}
	res, err := client.ConvertImages(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
	fmt.Println("done")
}

func CreateImage() {
	fmt.Println("create image")
	req := &imageservicepb.CreateImageRequest{}

	for {
		fmt.Print("image file path: ")
		scanner.Scan()
		imageFileName := scanner.Text()
		imageFilePath := fmt.Sprintf("./img/%s", imageFileName)

		bytes, err := os.ReadFile(imageFilePath)
		if err != nil {
			fmt.Println("file not found")
			continue
		}

		req.ImageFile = bytes
		req.ObjectName = imageFileName

		break
	}

	for {
		fmt.Print("resize width percent: ")
		scanner.Scan()
		input := scanner.Text()
		i, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("please enter number")
			continue
		}
		if i < 0 || i > 100 {
			fmt.Println("please enter number between 0 and 100")
			continue
		}

		req.ResizeWidthPercent = int32(i)

		break
	}

	for {
		fmt.Print("resize height percent: ")
		scanner.Scan()
		input := scanner.Text()
		i, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("please enter number")
			continue
		}
		if i < 0 || i > 100 {
			fmt.Println("please enter number between 0 and 100")
			continue
		}

		req.ResizeHeightPercent = int32(i)
		break
	}

	for {
		fmt.Print("encode format (1: JPEG or 2: PNG): ")
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "1":
			req.EncodeFormat = imageservicepb.EncodeFormat_JPEG
		case "2":
			req.EncodeFormat = imageservicepb.EncodeFormat_PNG
		default:
			fmt.Println("please enter JPEG or PNG")
			continue
		}

		break
	}

	res, err := client.CreateImage(context.Background(), req)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
	fmt.Println("done")
}
