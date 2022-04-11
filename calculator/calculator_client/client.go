package main

import (
	"context"
	"fmt"
	"log"

	"github.com/andiahmads/simple-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello im calculator client")

	//buat koneksi ke calculator Server
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)
	// DoUnary(c)

	//do client streaming
	DoClientStreaming(c)

}

func DoUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do unary RPC...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  80,
		SecondNumber: 40,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling RPC %v", err)
	}
	log.Fatalf("response from Greet %v", res.SumResult)

}

func DoClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do  a computeAverage client streaming.....")

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}
	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		fmt.Printf("Sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receipt response: %v", err)
	}

	fmt.Printf("the average is: %v", res.GetAverage())

}
