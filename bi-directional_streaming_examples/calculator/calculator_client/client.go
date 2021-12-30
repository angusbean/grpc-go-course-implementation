package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/angusbean/grpc-go-course/client_streamling_examples/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Client running\n")

	cc, err := grpc.Dial("localhost:50053", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	//doUnary(c)

	//doServerStreaming(c)

	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting Unary RPC...")
	req := &calculatorpb.SumRequest{
		Sum: &calculatorpb.Sum{
			IntOne: 5,
			IntTwo: 10,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Calculator RPC: %v", err)
	}
	log.Printf("Response from Calculator: %v", res.Result)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {

	fmt.Println("Starting to do a PrimeDecomposition Server Streaming RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 200,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happen: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
	}
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v\n", err)
	}
	for _, req := range requests {
		log.Printf("Sending request: %v\n", req.Number)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving response from ComputeAverage: %v\n", err)
	}
	fmt.Printf("Computed Average: %v\n", res.Result)
}
