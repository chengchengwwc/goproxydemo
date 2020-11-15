package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	pb "goproxywork/grpc_server_client/proto"
	"io"
	"log"
	"sync"
	"time"
)

const (
	timestampFormat = time.StampNano
	streamingCount  = 10
	AccessToken     = "eyJhbGciOiJIUzI1NiIsInR5cCI"
)

var addr = flag.String("addr", "localhost:8402", "the address to connect to")

func unaryCallWithMetadata(c pb.EchoClient, message string) {
	fmt.Println("--- unary --- \n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	r, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("failed to call unaryEcho :%v", err)
		return
	}
	fmt.Printf("response:%v \n", r.Message)
}

func serverStreamingWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("--- server streaming ---\n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.ServerStreamingEcho(ctx, &pb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf(("bad"))
		return
	}
	var rpcStatus error
	fmt.Printf("response:\n")
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf(" - %s\n", r.Message)
	}

	if rpcStatus != io.EOF {
		log.Printf("bad")
	}
}

func clientStreamWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("--- client streaming ---\n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.ClientStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("bad")
		return
	}

	for i := 0; i < streamingCount; i++ {
		if err := stream.Send(&pb.EchoRequest{Message: message}); err != nil {
			log.Fatalf("failed to send streaming: %v\n", err)
		}
	}

	r, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to CloseAndRecv: %v\n", err)
		return
	}
	fmt.Printf("response:%v\n", r.Message)
}

func bidirectionalWithMetadata(c pb.EchoClient, message string) {
	fmt.Printf("--- bidirectional ---\n")
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	md.Append("authorization", "Bearer "+AccessToken)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := c.BidirectionalStreamingEcho(ctx)
	if err != nil {
		log.Fatalf("bad")
		return
	}
	go func() {
		for i := 0; i < streamingCount; i++ {
			err = stream.Send(&pb.EchoRequest{Message: message})
			if err != nil {
				log.Fatalf("failed to send streaming: %v\n", err)
			}
		}
		stream.CloseSend()
	}()
	var rpcStatus error
	fmt.Printf("response:\n")
	for {
		r, err := stream.Recv()
		if err != nil {
			rpcStatus = err
			break
		}
		fmt.Printf(" - %s\n", r.Message)
	}
	if rpcStatus != io.EOF {
		log.Fatalf("failed to finish server streaming: %v", rpcStatus)
	}
}

func main() {
	flag.Parse()
	message := "hello weicheng"
	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn, err := grpc.Dial(*addr, grpc.WithInsecure())
			if err != nil {
				return
			}
			defer conn.Close()

			c := pb.NewEchoClient(conn)
			// step one
			unaryCallWithMetadata(c, message)
			time.Sleep(400 * time.Millisecond)

			serverStreamingWithMetadata(c, message)
			time.Sleep(1 * time.Second)

			clientStreamWithMetadata(c, message)
			time.Sleep(1 * time.Second)

			bidirectionalWithMetadata(c, message)
		}()
	}
	wg.Wait()
	time.Sleep(1 * time.Second)

}
