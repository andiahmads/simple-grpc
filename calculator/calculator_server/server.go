package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/andiahmads/simple-grpc/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Receipt Sun RPC:%v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SumResult: sum,
	}
	fmt.Printf("result:%v", res)
	return res, nil

}

func (*server) PrimeNumberDecompositon(req *calculatorpb.PrimeNumberDecompositonRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositonServer) error {
	fmt.Printf("Receipt PrimeNumberDecompositon RPC:%v", req)

	number := req.GetNumber()
	divisor := int64(2)
	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompositonResponse{
				PrimeFactor: int32(divisor),
			})

			number = number / divisor
		} else {
			divisor++
			fmt.Printf("DIVISOR has increased to %v\n", divisor)
		}
	}
	return nil

}

func main() {
	fmt.Println("running calculator server......")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed listening: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v \n", err)
	}

}
