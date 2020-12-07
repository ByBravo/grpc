package main

import (
	"context"
	"io"

	greetpb "github.com/ByBravo/grpc/grpc-course-udey/greet/greetpb"
	"github.com/ByBravo/grpc/grpc-course-udey/greet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var loggerf = log.LoggerJSON().WithField("package", "client")

func main() {

	log := loggerf.WithField("func", "main")
	log.Info("Hello im a client")

	tls := false
	opts := grpc.WithInsecure()
	if tls {
		certFile := "ssl/ca.crt" // Certificate Authority Trust certificate
		creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
		if sslErr != nil {
			log.Error("Error while loading CA trust certificate: " + sslErr.Error())
			return
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Error("could not connect: " + err.Error())
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)

	doUnary(c)
	doServerStreaming(c)
	// doClientStreaming(c)
	// doBiDiStreaming(c)

	// doUnaryWithDeadline(c, 5*time.Second) // should complete
	// doUnaryWithDeadline(c, 1*time.Second) // should timeout
}

func doUnary(c greetpb.GreetServiceClient) {
	log := loggerf.WithField("func", "doUnary")
	log.Info("Starting rpc call")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Byron",
			LastName:  "Bravo",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Error("error whilecalling Greet RPC: " + err.Error())
	}

	log.Info("Response from Greet " + res.Result)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	log := loggerf.WithField("func", "doServerStreaming")
	log.Info("Starting rpc call")
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Byron",
			LastName:  "Bravo",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Error("error whilecalling Greet RPC: " + err.Error())
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			log.Error("error whilecalling GreetMenyTimes : " + err.Error())

			break
		}

		if err != nil {
			log.Error("error whilecalling GreetMenyTimes RPC: " + err.Error())
		}

		log.Info("Response from GreetMenyTimes " + msg.GetResult())

	}

}
