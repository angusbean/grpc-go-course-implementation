package main

import (
	"context"
	"fmt"
	"log"

	"github.com/angusbean/grpc-go-course/uniary_examples/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Client running\n")

	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)

	doUnary(c)
}

// doUnary is a Uniary example of sending a GreetRequest in order to receive a GreetResponse from the server
func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Unary RPC...")
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
