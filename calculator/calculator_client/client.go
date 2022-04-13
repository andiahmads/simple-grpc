package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/andiahmads/simple-grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// DoClientStreaming(c)

	// DoBidiStreaming(c)

	//do error unary
	DoErrorUnary(c)

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

func DoBidiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("starting to do a FindMaximum Bidi streaming RPC....")

	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)

	}

	waitc := make(chan struct{})

	//send go routine

	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("sending number:%v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	//receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while reading server stream: %v", err)
				break
			}
			maximum := res.GetMaximum()
			fmt.Printf("Received a new maximum of...%v \n", maximum)
		}
		close(waitc)
	}()

	<-waitc

}

func DoErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a square root Unary RPC......")

	// case success
	// doErrorCall(c, 10)

	// case error
	doErrorCall(c, -2)

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {

	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})
	if err != nil {
		respErr, ok := status.FromError(err)

		if ok {
			//actual error from GRPC (user error)
			fmt.Printf("error message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("we probably sent a negative number!")
			}

		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Result of square root Unary RPC: %v", n, res.GetNumberRoot())
}
