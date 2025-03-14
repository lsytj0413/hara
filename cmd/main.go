// Copyright (c) 2024 The Songlin Yang Authors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

// Package main is the entrance of project
package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lsytj0413/golang-project-template/pb"
	"github.com/lsytj0413/golang-project-template/pkg/server"
	"github.com/lsytj0413/golang-project-template/pkg/utils/version"
)

// curl -v -H "Request-ID: 123" http://127.0.0.1:8081/v1/hello -d '{"Name": "lsytj0413"}'
// *   Trying 127.0.0.1...
// * TCP_NODELAY set
// * Connected to 127.0.0.1 (127.0.0.1) port 8081 (#0)
// > POST /v1/hello HTTP/1.1
// > Host: 127.0.0.1:8081
// > User-Agent: curl/7.64.1
// > Accept: */*
// > Request-ID: 123
// > Content-Length: 21
// > Content-Type: application/x-www-form-urlencoded
// >
// * upload completely sent off: 21 out of 21 bytes
// < HTTP/1.1 200 OK
// < Content-Type: application/json
// < Request-Id: r_123
// < Date: Thu, 30 Mar 2023 02:11:57 GMT
// < Content-Length: 74
// <
// * Connection #0 to host 127.0.0.1 left intact
// {"Message":"Hello lsytj0413!","CurrentTime":"2023-03-30T02:11:57.619397Z"}* Closing connection 0.
func main() {
	fmt.Printf("%s\n", version.Get().Pretty())

	//nolint
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterHelloServiceServer(grpcServer, server.NewHelloServer())
	//nolint
	go grpcServer.Serve(lis)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(s string) (string, bool) {
			fmt.Printf("receieve request header: %s\n", s)
			return s, true
		}),
		runtime.WithOutgoingHeaderMatcher(func(s string) (string, bool) {
			fmt.Printf("receieve response header: %s\n", s)
			return s, true
		}),
	)
	gwopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, "0.0.0.0:8080", gwopts)
	if err != nil {
		panic(err)
	}

	//nolint
	http.ListenAndServe("0.0.0.0:8081", mux)
}
