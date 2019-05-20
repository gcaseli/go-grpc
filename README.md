# go-grpc
Project about gRPC protocol

###  Whats is gRPC? ### 
  It’s a free open source created by google to define request and response for RPC (Remote Procedure Calls)
  Build on top of HTTP/2, low latency, support streaming

### What is Protocol Buffer? ### 
Are a way of encoding structured  data in a efficient yet extensible format

### Why use gRPC? ### 
  Easy code definition
  Use HTTP/2 as transport mechanism
  Support for streaming API for maximum performance
  Is API oriented, instead of Resource Oriented like REST
  
### GO DEPENDENCIES ### 

```
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```


### What is server streaming API? ### 
 - It’s a new kind API enabled thanks to HTTP/2
 - The client will send one message to the server and will receive many responses from the server, possibly any infinite number
 - Streaming Server are well suited for:
 - When the server needs to send a lot of data (big data)
 - When the server needs to PUSH data to the client without having the client request for more (think live feed, chat)
 - In gRPC Server Streaming Calls are defined using the keyword “stream”

### What is client streaming API? ### 
 - It’s a new kind API enabled thanks to HTTP/2
 - The client will send many message to the server and will receive one response from the server (at any time)
 - Streaming client are well suited for:
 - When the client needs to send a lot of data
 - When the server processing is expensive and should happen as the client send datas.
 - When the client needs to PUSH data to the server without really expecting a response

### What is Bi Directional Streaming API? ### 
 - It’s a new kind API enabled thanks to HTTP/2
 - The client will send many message to the server and will receive many response from the server
 - The number of requests and response does not have to match
 - When the client and the server needs to send a lot of data asynchronously
 - Long running connections
 - “Chat” protocol

