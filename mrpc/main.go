package main

import (
	"mp/mrpc/mproto"
	"golang.org/x/net/context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"io"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/metadata"
	"net/http"
	"golang.org/x/net/trace"
)

type Svr struct {

}

func (s *Svr)ReqSvr(ctx context.Context, req *pb.Req) (*pb.Ack, error) {
	fmt.Printf("A=%v, B=%v\n", req.A, req.B)
	ack := &pb.Ack{
		V:fmt.Sprintf("A=^v", req.A),
	}

	return ack, nil
}

func (s *Svr)StreamSvr(stream pb.Test_StreamSvrServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("err =%v", err)
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Println("req A=%v, B=%v\n", in.A, in.B)
		ack := &pb.Ack{
			V:fmt.Sprintf("A=%v, B=%v", in.A, in.B),
		}

		stream.Send(ack)
	}

	return nil
}

func newServer() *Svr {
	s := &Svr{}
	return s
}

// auth 验证Token
func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}

	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}

	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return status.Errorf(codes.Unauthenticated, "Token认证信息无效: appid=%s, appkey=%s", appid, appkey)
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 8080))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTestServer(grpcServer, newServer())

	go startTrace()

	grpcServer.Serve(lis)
}

func startTrace() {
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		fmt.Printf("trace authrequest ok\n")
		return true, true
	}
	go http.ListenAndServe(":50051", nil)
	fmt.Printf( "%s", "Trace listen on 50051\n")
}
