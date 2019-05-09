package main

import (
	"context"
	"fmt"
	"go-grpc/unary/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I'm the client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do unary rpc...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Guilherme",
			LastName:  "Caseli",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("err when calling gRPC: %v", err)
	}

	log.Printf("Reponse from greet %v", res.Result)
}
