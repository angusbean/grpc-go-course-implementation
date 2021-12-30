# grpc-go-course-implementation

Completed gRPC exercises following the Udemy Course (https://www.udemy.com/course/grpc-golang/)

Examples of all 4 gRPC transports detailed:
    1) Unary
    2) Server Streaming
    3) Client Streaming
    4) Bi-directional Streaming

Issue with Protoc generating go-gen-go-grpc, ensure the installed plugin cna be found in hte $PATH add it with:
export PATH=$PATH:$(go env GOPATH)/bin

Ensure both the '\*_grpc.pb.go' and the '\*.pb.go' files have been created and up to date by running the commands in the generate.sh file

gRPC Error handling documentation
https://grpc.io/docs/guides/error/

Handy Guide for gRPC Errors
https://avi.im/grpc-errors/

Server

To send an error, return status.Errorf with error message and code:

status.Errorf(<grpc error code>, <error message>)

Example:

return status.Errorf(codes.InvalidArgument, "Ouch!")

Client

To handle the error, check error returned from gRPC call:

_, err := client.GRPCMethod(...)
statusCode := status.FromError(err)
