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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/durd07/grpc-test/tra"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedTraServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) Nodes(ctx context.Context, in *pb.TraRequest) (*pb.TraResponse, error) {
	log.Printf("Received: %v", in.Fqdn)
	return &pb.TraResponse{Fqdn: in.Fqdn, Nodes: []*pb.Node{
			&pb.Node{NodeId: "1", Ip: "192.168.0.1", SipPort: 5060, Weight: 50},
			&pb.Node{NodeId: "2", Ip: "192.168.0.2", SipPort: 5060, Weight: 50},
	}}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTraServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
