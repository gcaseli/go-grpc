package main

import (
	"context"
	"fmt"
	"go-grpc/blog/blogpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	fmt.Println("Creating blog")

	blog := &blogpb.Blog{
		AuthorId: "Guilherme Caseli",
		Title:    "My fisrt blog",
		Content:  "Anything",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{
		Blog: blog,
	})

	if err != nil {
		log.Fatalf("error when creating blog: %v", err)
	}

	fmt.Printf("Blog was created: %v", res)

}
