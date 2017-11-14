package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	pb "myapp/models/grpc/service/echo"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "192.168.34.134:50051"
)

func echoTimeService(c pb.EchoClient) {

	stream, err := c.EchoTime(context.Background(), &pb.Request{Message: "time stream"})

	if err != nil {
		log.Fatalf("could not greet: %v\n", err)
	}

	for {
		time, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("\nEcho time err: %v\n", err)
		}
		fmt.Printf("\ntime is: %s\n", time.Message)
	}
}

func echoService(c pb.EchoClient) {
	scan := bufio.NewReader(os.Stdin)
	stream, err := c.EchoHello(context.Background())
	if err != nil {
		log.Printf("\nEcho stream error: %v\n", err)
	}

	waitc := make(chan struct{})

	go func() {
		for {
			info, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("\nEcho service error:%v", err)
			}
			fmt.Printf("Info: %s\n", info)
		}
	}()

	for {
		mess, isProfix, err := scan.ReadLine()
		if err != nil {
			log.Fatalf("\nRead message error:%v", err)
		} else if isProfix {
			log.Fatalf("\n Your message out of size!")
		}
		message := string(mess[:])
		if err := stream.Send(&pb.Request{Message: message}); err != nil {
			log.Printf("Failed to send a note: %v\n", err)
		}
	}

}

func main() {
	// Set connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Not connect: %v\n", err)
	}
	defer conn.Close()
	c := pb.NewEchoClient(conn)
	go echoTimeService(c)
	echoService(c)

}
