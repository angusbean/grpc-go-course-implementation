package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/angusbean/grpc-go-course/complete_examples/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Client is running...")

	tls := false // if tls is true SSL will be enabled, if false will run as unsecure
	opts := grpc.WithInsecure()
	if tls {
		certFile := "../ssl/ca.crt" //Cert Auth Trust Certification
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50052", opts)

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)

	//doServerStreaming(c)

	//doClientStreaming(c)

	//doBiDirectionalStreaming(c)

	//doUnaryWithDeadline(c, 1*time.Second)

	//doUnaryWithDeadline(c, 10*time.Second)

}

func doUnary(c greetpb.GreetServiceClient) {

	fmt.Println("Starting Unary RPC...\n")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Angus",
			LastName:  "Bean",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Printf("Starting Server Streaming RPC...\n")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Angus",
			LastName:  "Bean",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Angus",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bill",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tony",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tim",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Phil",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}
	// we iterate over our slice and send each message individually
	for _, req := range requests {
		log.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a bi-directional Streaming RPC...")

	// Create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating steam: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Angus",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bill",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tony",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tim",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Phil",
			},
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
		// function to receive many messages
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

func doUnaryWithDeadline(c greetpb.GreetServiceClient, timeout time.Duration) {

	fmt.Println("Starting UnaryWithDeadline RPC...")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Angus",
			LastName:  "Bean",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)

		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Timeout expired, Deadline was exceeded")
			} else {
				fmt.Printf("Unexpected error: %v\n", statusErr)
			}
		} else {
			log.Fatalf("Error calling GreetWithDeadline: %v", err)
		}
		return
	}
	log.Printf("Response from GreetWithDeadline: %v", res.Result)
}
