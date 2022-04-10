package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doServerStreaming(c)

	//client streaming
	doClientStreaming(c)

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

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("starting client streaming........")

	request := []*greetpb.LongGreateRequest{
		&greetpb.LongGreateRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng",
			},
		},
		&greetpb.LongGreateRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 1",
			},
		},

		&greetpb.LongGreateRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 2",
			},
		},

		&greetpb.LongGreateRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 3",
			},
		},
		&greetpb.LongGreateRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 4",
			},
		},

		&greetpb.LongGreateRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 5",
			},
		},
	}

	stream, err := c.LongGreate(context.Background())
	if err != nil {
		log.Fatalf("error while reading stream client streaming:%v", err)
	}

	for _, req := range request {
		fmt.Printf("sending req:%v \n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)

	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receipt stream client streaming:%v", err)
	}

	fmt.Printf("LongGreet response:%v", res)

}
