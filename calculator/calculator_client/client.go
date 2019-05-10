package main

import (
	"context"
	"fmt"
	"go-grpc/calculator/calculatorpb"
	"io"
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

	doServerStreaming(c)
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

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do PrimeNumberDecomposition Streaming rpc...")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 3243243243234,
	}

	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("err when calling PrimeNumberDecomposition gRPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("somenthing happened: %v", err)
		}
		fmt.Println(res.GetPrimeFactor())
	}

}
