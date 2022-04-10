package main

import (
	"context"
	"fmt"
	"io"
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

	// unary api
	// DoUnary(c)

	//streaming server
	DoStreamingServer(c)

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
	log.Fatalf("response from Greet %v", res.GetSumResult())

}

func DoStreamingServer(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do Streaming RPC...")

	req := &calculatorpb.PrimeNumberDecompositonRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecompositon(context.Background(), req)
	if err != nil {
		if err != nil {
			log.Fatalf("error while calling RPC %v", err)
		}
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream:%v", err)
		}
		log.Println(msg.PrimeFactor)
	}

}
