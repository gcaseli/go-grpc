package main

import (
	"context"
	"fmt"
	"go-grpc/blog/blogpb"
	"io"
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

	blogID := res.GetBlog().GetId()

	fmt.Printf("Blog was created: %v", res)

	fmt.Println("Reading the blog")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "5d0a6faa3b05825269a2d56f",
	})

	if err2 != nil {
		fmt.Printf("Erro when reading blog, %v\n", err2)
	}

	readBlogResp, readBlogErr := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: blogID,
	})

	if readBlogErr != nil {
		fmt.Printf("Erro when reading blog, %v", err2)
	}

	fmt.Printf("blog is, %v\n", readBlogResp)

	fmt.Println("Updating the blog")

	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Guilherme Sampaio Caseli",
		Title:    "My second blog UPDATE",
		Content:  "Anything mooooorrree UPDATE",
	}

	updRest, updErr := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: newBlog,
	})

	if updErr != nil {
		fmt.Printf("Erro when updating blog, %v", updErr)
	}

	fmt.Printf("blog was updated , %v\n", updRest)

	fmt.Println("Deleting the blog with id: " + blogID)

	delRes, delErr := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: blogID,
	})

	if delErr != nil {
		fmt.Printf("Erro when deleting blog, %v", delErr)
	}

	fmt.Printf("blog was deleted , %v\n", delRes)

	fmt.Println("========================")

	fmt.Println("List Blog data")

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})

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
		fmt.Printf("Blog returned is: %v\n ", res.GetBlog())
	}

}
