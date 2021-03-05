/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/durd07/grpc-test/tra"
	"google.golang.org/grpc"
)

const (
	address     = "felixdu.hz.dynamic.nsn-net.net:50053"
	defaultFqdn = "tafe.default.svc.cluster.local"
)

//func main() {
//	// Set up a connection to the server.
//	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
//	if err != nil {
//		log.Fatalf("did not connect: %v", err)
//	}
//	defer conn.Close()
//	c := pb.NewTraServiceClient(conn)
//
//	// Contact the server and print out its response.
//	fqdn := defaultFqdn
//	if len(os.Args) > 1 {
//		fqdn = os.Args[1]
//	}
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	r, err := c.Nodes(ctx, &pb.TraRequest{Fqdn: fqdn})
//	if err != nil {
//		log.Fatalf("could not query node : %v", err)
//	}
//	log.Printf("get node: %v", r)
//}

func recvNotification(stream pb.TraService_SubscribeClient) {
	for {
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("failed to recv %v", err)
		}

		if err == io.EOF {
			break
		}

		log.Println(resp, err)
	}
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTraServiceClient(conn)

	// Contact the server and print out its response.
	fqdn := defaultFqdn
	if len(os.Args) > 1 {
		fqdn = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.Subscribe(&pb.TraRequest{Fqdn: fqdn})
	if err != nil {
		log.Fatalf("could not query node : %v", err)
	}

	recvNotification(stream)
	stream.CloseSend()
}
