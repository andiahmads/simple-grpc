package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/andiahmads/simple-grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doClientStreaming(c)

	//bidi streaming
	// doBidiStreaming(c)

	//with deadline
	// DoUnaryWithDeadline(c, 5*time.Second) //complate
	DoUnaryWithDeadline(c, 1*time.Second) // timeout

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

func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do bidi streaming RPC.....")

	//create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while creating stream client:%v", err)
		return
	}

	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 1",
			},
		},

		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 2",
			},
		},

		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 3",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 4",
			},
		},

		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Andi ganteng 5",
			},
		},
	}

	waitc := make(chan struct{})

	// send a bunch of message to the client (go routine)
	go func() {
		//function to send a bunc of message

		for _, req := range request {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()

	}()

	//receive a bunch of message from the client (go routine)
	go func() {

		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)

	}()

	// block until everything is done
	<-waitc

}

func DoUnaryWithDeadline(c greetpb.GreetServiceClient, timeOut time.Duration) {
	fmt.Println("Starting to do unary with deadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "andiahmad",
			LastName:  "saputra",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout was hit! Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error:%v", err)
			}
		} else {

			log.Fatalf("error while calling RPC %v", err)
		}
	}
	log.Fatalf("response from Greet %v", res.Result)

}
