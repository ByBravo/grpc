package main

import (
	"context"
	"fmt"
	"io"

	//"log"
	"net"
	"strconv"
	"time"

	greetpb "github.com/ByBravo/grpc/grpc-course-udey/greet/greetpb"
	"github.com/ByBravo/grpc/grpc-course-udey/greet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

var loggerf = log.LoggerJSON().WithField("package", "server")

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	log := loggerf.WithField("func", "Greet")
	log.Info("Starting")
	firstName := req.GetGreeting().GetFirstName()

	if firstName == "" {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a empty firstname "),
		)
	}

	result := "Hello " + firstName

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log := loggerf.WithField("func", "GreetManyTimes")
	log.Info("Starting")
	//
	stringRequest := fmt.Sprintf("%#v", req)
	log.Debug("Function was invoked with : " + stringRequest)

	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreetManytimesResponse{
			Result: result,
		}
		//PAY attention here
		log.Info("stream.Send")
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	log := loggerf.WithField("func", "LongGreet")
	log.Info("Starting")

	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {

	log := loggerf.WithField("func", "GreetEveryone")
	log.Info("Starting")
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	log := loggerf.WithField("func", "GreetWithDeadline")
	log.Info("Starting")
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	log := loggerf.WithField("func", "main")
	log.Info("Starting")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Error("Failed to listen: " + err.Error())
		return
	}

	opts := []grpc.ServerOption{}
	tls := false
	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
		if sslErr != nil {
			log.Error("Failed loading certificates : " + sslErr.Error())
			return
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Error("Failed to serve: " + err.Error())
	}
}
