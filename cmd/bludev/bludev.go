package bludev

import (
	"context"
	"fmt"
	"io"
	"log"

	pb "github.com/flood-io/cli/proto"
	"google.golang.org/grpc"
)

type BLUDev struct {
}

type testServer struct {
}

func (t *testServer) Run(*pb.TestRequest, pb.Test_RunServer) error {
	return nil
}

var _ pb.TestServer = (*testServer)(nil)

func (b *BLUDev) Run(scriptFile string) {
	fmt.Println("running dev-blu")
	fmt.Printf("scriptFile = %+v\n", scriptFile)

	serverAddr := "localhost:50051"

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTestClient(conn)

	fmt.Printf("client = %+v\n", client)

	test := &pb.TestRequest{
		ClientID: "123",
		Uuid:     "456",
		Script:   []byte("hey its a script"),
	}

	fmt.Println("streaming")
	stream, err := client.Run(context.Background(), test)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.Run(_) = _, %v", client, err)
		}
		fmt.Printf("result = %+v\n", result)
		fmt.Printf("result = %T\n", result)
		log.Println("result", result.String())
	}
}
