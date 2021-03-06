package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/angusbean/grpc-go-course/client_streamling_examples/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage invoked with a streaming request")

	totalValue := float64(0)
	recvCount := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			result := totalValue / float64(recvCount)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while receiving client stream: %v", err)
		}
		number := float64(req.GetNumber())
		totalValue = totalValue + number
		recvCount++
	}
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
