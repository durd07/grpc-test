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

	"google.golang.org/grpc"
	pb "github.com/durd07/grpc-test/tra"
)

const (
	address     = "localhost:50051"
	defaultFqdn = "tafe.default.cluster.local"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTraClient(conn)

	// Contact the server and print out its response.
	fqdn := defaultFqdn
	if len(os.Args) > 1 {
		fqdn = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Nodes(ctx, &pb.TraRequest{fqdn: fqdn})
	if err != nil {
		log.Fatalf("could not query node : %v", err)
	}
	log.Printf("get node: %s %v", r.GetFqdn(), r.GetNodes())
}
