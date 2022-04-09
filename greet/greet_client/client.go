package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/andiahmads/simple-grpc/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello im client")

	//buat connection ke server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	// log.Fatalf("created Client: %f", c)

	//unary API
	// doUnary(c)

	//server streaming
	doServerStreaming(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do unary RPC...")

	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "andiahmad",
			LastName:  "saputra",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling RPC %v", err)
	}
	log.Fatalf("response from Greet %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting server streaming........")

	req := &greetpb.GreateManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "andiahmad",
			LastName:  "saputra",
		},
	}

	//dependency injection
	resStream, err := c.GreateManyTime(context.Background(), req)
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
		log.Printf("Response from greatmanytime: %v", msg.GetResult())
	}

}
