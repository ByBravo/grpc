package main

import (
	"context"
	"io"
	"time"

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

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDiStreaming(c)

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

func doClientStreaming(c greetpb.GreetServiceClient) {
	log := loggerf.WithField("func", "doClientStreaming")
	log.Info("Starting rpc call")
	request := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Byron",
				LastName:  "Bravo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Natalia",
				LastName:  "Espildora",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lorna",
				LastName:  "Bravo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ivar",
				LastName:  "Bravo",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Error("error while calling LongGreet RPC: " + err.Error())
	}

	for _, req := range request {
		log.Info("Sending req :" + req.Greeting.FirstName)
		stream.Send(req)

		time.Sleep(100 * time.Millisecond)

	}

	res, err := stream.CloseAndRecv()

	if nil != err {
		log.Error("error while receiving response from LongGreet : " + err.Error())
	}

	log.Info("LongGreet Response : " + res.Result)

}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	log := loggerf.WithField("func", "doBiDiStreaming")
	log.Info("Starting rpc call")

	request := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Byron",
				LastName:  "Bravo",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Natalia",
				LastName:  "Espildora",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lorna",
				LastName:  "Bravo",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Ivar",
				LastName:  "Bravo",
			},
		},
	}
	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Error("error while calling creating stream : " + err.Error())
	}

	waitc := make(chan struct{})
	// we send a bunch of msg to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range request {
			log.Info("Sending message :" + req.Greeting.FirstName)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)

		}
		stream.CloseSend()

	}()

	// we receive a bunch of msg from the client (go routine)
	go func() {

		//LIKE A CHAT
		for {

			res, err := stream.Recv()
			if err == io.EOF {
				log.Error("error while receiving GreetEveryone : " + err.Error())
				//ENABLE THIS TO RECEIVE FIRST ITEM
				//close(waitc)
				break
			}

			if err != nil {
				log.Error("error whilecalling GreetEveryone RPC: " + err.Error())
				//
				break
			}

			log.Info("GreetEveryone Response : " + res.Result)
		}
		//ENABLE THIS TO RECEIVE ALL RESPONSES
		close(waitc)
	}()

	// block until everything is done
	<-waitc

}
