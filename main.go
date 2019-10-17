/*
MIT License
Copyright (c) 2019 kurosawa
Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"context"
	"log"
	"time"
	"io/ioutil"
	"google.golang.org/grpc"
	"flag"
	"net"
	"fmt"
	"path"
	pb "upload/proto"
)

const (
	port = ":50051"
)

type server struct{}

func request(address , file string) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUploadClient(conn)

	if file == "" {
		log.Fatal("Please Specify Sending File")
	}

	data, err := ioutil.ReadFile(file)
    if err != nil {
        log.Fatal(err)
    }
	basefile := path.Base(file)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Fileup(ctx, &pb.Request{Data: data, Name: basefile})
	if err != nil {
		log.Fatalf("could not Sending File: %v", err)
	}
	log.Printf(": %s", r.GetMessage())
}

func (s *server) Fileup(ctx context.Context, in *pb.Request) (*pb.Reply, error) {
	log.Printf("Received: %v", in.GetName())
	err := ioutil.WriteFile(in.GetName(), in.GetData(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return &pb.Reply{Message: "Upload Done!!"}, nil
}

func reply() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Println("waiting file..")
	s := grpc.NewServer()
	pb.RegisterUploadServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


func main() {
	ip := flag.String("i", "", "IP address")
	file := flag.String("f", "", "filename")
	flag.Parse()
	address := *ip + port
	if *ip == "" && *file == "" {
		reply()
	} else if *ip != "" && *file == "" {
		log.Fatal("Please Set File Name")
	} else if *ip == "" && *file != "" {
		log.Fatal("Please Set IP Adress")
	} else {
		request(address, *file)
	}
}
