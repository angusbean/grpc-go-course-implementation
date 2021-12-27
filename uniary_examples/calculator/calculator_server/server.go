package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/angusbean/grpc-go-course/uniary_examples/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstInt := req.Sum.IntOne
	secondInt := req.Sum.IntTwo

	result := firstInt + secondInt

	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil

}

func main() {
	fmt.Println("Server is running")

	lis, err := net.Listen("tcp", "0.0.0.0:50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
