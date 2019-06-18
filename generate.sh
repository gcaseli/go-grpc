protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.
protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.
protoc blog/blogpb/blog.proto --go_out=plugins=grpc:.

//driver for mongodb
go get github.com/mongodb/mongo-go-driver/mongo
go get go.mongodb.org/mongo-driver/mongo/options

--mogno
cd /home/ent_gcaseli/Downloads/mongodb-linux-x86_64-ubuntu1604-4.0.8
bin/mongod --dbpath data/db