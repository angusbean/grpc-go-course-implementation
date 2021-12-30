package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/angusbean/grpc-go-course/bi-directional_streaming_examples/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum invoked with a bi-directional stream request")

	maxNumValue := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading from client stream: %v\n", err)
		}

		newNumValue := req.GetNumber()

		if newNumValue > maxNumValue {
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				Result: newNumValue,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
			maxNumValue = newNumValue
		}
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
