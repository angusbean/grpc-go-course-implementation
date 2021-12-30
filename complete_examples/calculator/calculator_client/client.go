package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/angusbean/grpc-go-course/complete_examples/calculator/calculatorpb"
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

	//doClientStreaming(c)

	doBiDirectionalStreaming(c)
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

func doBiDirectionalStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a bi-directional streaming RPC...")

	// Create a stream by invoking the client
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 1,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 10,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 5,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 10,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 1000,
		},
	}

	waitChannel := make(chan struct{})

	// Send messages to the client
	go func() {
		// function to send many messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive messages from the client
	go func() {
		// function to receive messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v\n", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitChannel)
	}()

	// Block until everything is finished
	<-waitChannel
}
