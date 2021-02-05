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

	"io/ioutil"
	"fmt"
	"net/http"
	"google.golang.org/grpc"
	"encoding/json"
	pb "github.com/durd07/grpc-test/tra"
)

const (
	port = ":50053"
)

type NodeData struct {
	Fqdn string `json:"fqdn"`
	Node_id string `json:"node_id"`
	Ip string `json:"ip"`
	Sip_port uint32 `json:"sip_port"`
	Weight uint32 `json:"weight"`
}

var (
	data = map[string] *pb.TraResponse {
		"tafe.svc.cluster.local" : &pb.TraResponse{
			Fqdn: "tafe.svc.cluster.local",
			Nodes: []*pb.Node{
				&pb.Node{NodeId: "1", Ip: "192.168.0.1", SipPort: 5060, Weight: 50},
				&pb.Node{NodeId: "2", Ip: "192.168.0.2", SipPort: 5060, Weight: 50},
			},
		},

		"scfe.svc.cluster.local" : &pb.TraResponse{
			Fqdn: "scfe.svc.cluster.local",
			Nodes: []*pb.Node{
				&pb.Node{NodeId: "3", Ip: "192.168.0.3", SipPort: 5060, Weight: 50},
				&pb.Node{NodeId: "4", Ip: "192.168.0.4", SipPort: 5060, Weight: 50},
			},
		},
	}
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedTraServiceServer
}


// SayHello implements helloworld.GreeterServer
func (s *server) Nodes(ctx context.Context, in *pb.TraRequest) (*pb.TraResponse, error) {
	log.Printf("Received: %v", in.Fqdn)
	var response = data[in.Fqdn]
	return response, nil
//	return &pb.TraResponse{Fqdn: in.Fqdn, Nodes: []*pb.Node{
//			&pb.Node{NodeId: "1", Ip: "192.168.0.1", SipPort: 5060, Weight: 50},
//			&pb.Node{NodeId: "2", Ip: "192.168.0.2", SipPort: 5060, Weight: 50},
//	}}, nil
}

func grpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTraServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func updateNodes(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
	fmt.Println("XXX")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Read failed:", err)
	}
	defer r.Body.Close()

	var update_data []NodeData
	err = json.Unmarshal(b, &update_data)
	if err != nil {
		log.Println("json format error:", err)
	}
	log.Println("update_data:", update_data)

	data = make(map[string]*pb.TraResponse)
	for _, v := range update_data {
		if _, ok := data[v.Fqdn]; ok {
			data[v.Fqdn].Nodes = append(data[v.Fqdn].Nodes, &pb.Node{NodeId: v.Node_id, Ip: v.Ip, SipPort: v.Sip_port, Weight: v.Weight})
		} else {
		data[v.Fqdn] = &pb.TraResponse{
			Fqdn: v.Fqdn,
			Nodes: []*pb.Node{
				&pb.Node{NodeId: v.Node_id, Ip: v.Ip, SipPort: v.Sip_port, Weight: v.Weight},
			},
		}
	}
	}

	for k, v := range data {
		fmt.Println(k, v.Nodes)
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
        return
    } else if r.Method == "GET" {
	    var tmp []NodeData
	for _, v := range data {
		for _, n := range v.Nodes {
			tmp = append(tmp, NodeData{v.Fqdn, n.NodeId, n.Ip, n.SipPort, n.Weight})
		}
	}
	fmt.Println(tmp)
	   b ,err := json.Marshal(tmp)
    if err != nil {
        log.Println("json format error:", err)
        return
    }

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
    w.Write(b)
    } else {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
        return
    }
}

func httpServer() {
    http.HandleFunc("/update", updateNodes)
    http.ListenAndServe(":50052", nil)
}

func main() {
	go grpcServer()
	httpServer()
}
