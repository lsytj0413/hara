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

// Package server provide the HelloServer implement
package server

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/lsytj0413/golang-project-template/pb"
	"github.com/lsytj0413/golang-project-template/pkg/utils"
)

// Reponser ...
type Reponser[T any] interface {
	Message(T) T
}

//nolint
type HelloServer struct {
	pb.UnimplementedHelloServiceServer
}

//nolint
func NewHelloServer() *HelloServer {
	return &HelloServer{}
}

//nolint
func (s *HelloServer) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		//nolint
		return nil, fmt.Errorf("MUST specify incoming metadata")
	}

	requestID := func() string {
		ids := md.Get("Request-Id")
		if len(ids) > 0 {
			return ids[0]
		}

		return ""
	}()

	if err := grpc.SendHeader(ctx, metadata.New(map[string]string{
		"Request-Id": "r_" + requestID,
	})); err != nil {
		return nil, fmt.Errorf("cannot send header for response, %w", err)
	}

	return &pb.HelloResponse{
		Message:     utils.GenerateResponseMessage(in.Name),
		CurrentTime: timestamppb.Now(),
	}, nil
}
