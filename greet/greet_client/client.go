package main

import (
	"context"
	"fmt"
	"go-grpc/greet/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("I'm the client")

	opts := grpc.WithInsecure()
	tls := true
	if tls {
		certFile := "ssl/ca.crt" // Certificate Authorization Trust
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Fatalf("error when loading certification, %v", sslErr)
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	//doBiDiStreaming(c)
	//doUnaryWithDeadline(c, 5*time.Second) // timeout
	//doUnaryWithDeadline(c, 1*time.Second) // complete
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

func doServerStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do Server Streaming RPC...")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Guilherme",
			LastName:  "Caseli",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("err when calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("err while reading stream ,%v", err)
		}
		log.Printf("response from GreetManyTimes: %v", msg.GetResult())
	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Guilherme",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Marjory",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Matheus",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Gustavo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "MundoDoce",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())

	if err != nil {
		log.Fatalf("err while calling LongGreet ,%v", err)
	}

	// we interate over our slace and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending the request: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("err while response from LongGreet ,%v", err)
	}
	fmt.Printf("LongGreet Response %v", response)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {

	fmt.Println("Starting to do BiDi Streaming RPC...")

	stream, err := c.GreetEveryone(context.Background())

	if err != nil {
		log.Fatalf("err while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Guilherme",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Marjory",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Matheus",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Gustavo",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "MundoDoce",
			},
		},
	}

	waitChannel := make(chan struct{})

	// send a bunch of messages to the server
	go func() {
		//function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("sending messages: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// receive a bunch of messages from the server
	go func() {
		//function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("err while receiving stream: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitChannel)
	}()

	//block until everything is done
	<-waitChannel
}

func doUnaryWithDeadline(c greetpb.GreetServiceClient, seconds time.Duration) {
	fmt.Println("Starting to do doUnaryWithDeadline rpc...")
	req := &greetpb.GreetDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Guilherme",
			LastName:  "Caseli",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), seconds)

	defer cancel()

	res, err := c.GreetDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Deadline was exceeded")
			} else {
				fmt.Printf("unexpected error: %v", statusErr)
			}
		} else {
			log.Fatalf("err when calling gRPC: %v", err)
		}
		return
	}

	log.Printf("Reponse from greet %v", res.Result)
}
