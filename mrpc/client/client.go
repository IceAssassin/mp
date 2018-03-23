package main

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
	"mp/mrpc/mproto"
	"fmt"
	"time"
	"io"
)

type client struct {

}

func do (cli  pb.TestClient, A string, B int64) {
	in := & pb.Req {
		A: A,
		B: B,
	}
	fmt.Printf("do (%v, %v)\n", in.A, in.B)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()



	res, err := cli.ReqSvr(ctx, in)
	if err != nil {
		fmt.Printf("err = %v\n", err)
		return
	}

	fmt.Sprintf("res = %v\n", res.V)
}

func dostream(cli pb.TestClient, A string, B int64) {
	in := & pb.Req {
		A: A,
		B: B,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := cli.StreamSvr(ctx)
	if err != nil {
		fmt.Println(" error = %v\n", err)
		return
	}

	//waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				time.Sleep(time.Second)
				continue
			}
			if err != nil {
				fmt.Printf("Failed to receive a note : %v\n", err)
				time.Sleep(time.Second)
				continue
			}
			fmt.Printf("A(%v:%v) Got message V=%v\n", A, B, in.V)
		}
	}()

	for i:=0; i<30; i++ {
		if err := stream.Send(in); err != nil {
			fmt.Printf("Failed to send a note: %v\n", err)
		}

		time.Sleep(time.Second)
	}

	stream.CloseSend()
	fmt.Printf("write close send\n")
}

// customCredential 自定义认证
type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return true
}


func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) //grpc.WithPerRPCCredentials(new(customCredential)))
	conn, err := grpc.Dial("127.0.0.1:8080", opts...)
	if err != nil {
		fmt.Println("fail to dial: %v", err)
		return
	}



	defer conn.Close()
	client := pb.NewTestClient(conn)

	//go dostream(client, "one", 10)
	//dostream(client, "tow", 10)
	//
	go func() {
		for {
			do(client, "one", 10)
			time.Sleep(time.Second * 1)
		}
	}()

	//go func() {
	//	for {
	//		do(client, "tow", 11)
	//		time.Sleep(time.Second * 1)
	//	}
	//}()

	select {

	}
}
