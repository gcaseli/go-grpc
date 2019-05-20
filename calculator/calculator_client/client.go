package main

import (
	"context"
	"fmt"
	"go-grpc/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Calculator Client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)
	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	doErrorUnary(c)
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
	fmt.Println("Starting to do PrimeNumberDecomposition Server Streaming rpc...")

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

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do ComputeAverage Client Streaming rpc...")

	stream, err := c.ComputeAverage(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	numbers := []int32{3, 5, 7, 54, 34}

	for _, number := range numbers {
		fmt.Printf("sending number: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while receving response: %v", err)
	}

	fmt.Printf("The average is: %v\n", res.GetAverage())
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do FindMaximun BiDi Streaming rpc...")

	stream, err := c.FindMaximun(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximun: %v", err)
	}

	waitc := make(chan struct{})

	go func() {
		numbers := []int32{4, 7, 2, 19, 4, 6, 32}
		for _, number := range numbers {
			fmt.Printf("Sending number: %v\n", number)
			stream.Send(&calculatorpb.FindMaximunResquest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading server streaming: %v", err)
				break
			}
			maximun := res.GetMaximun()
			fmt.Printf("Received a new maximun of....%v\n", maximun)
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do SquareRoot Streaming rpc...")

	//call ok
	doErrorCall(c, 10)

	// call error
	doErrorCall(c, -10)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC
			fmt.Printf("Return from the server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a  negative number!")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoor: %v", err)
			return
		}
	}

	fmt.Printf("Resulf of SquareRoot of %v: %v\n", number, res.GetNumberRoot())
}
