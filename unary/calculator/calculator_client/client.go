package main

import (
	"context"
	"fmt"
	"go-grpc/unary/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do SUM unary rpc...")

	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 4,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("err when calling SUM gRPC: %v", err)
	}

	log.Printf("Reponse from SUM %v", res.SumResult)
}
