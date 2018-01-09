package bludev

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

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

func (b *BLUDev) Run(scriptFile string) (err error) {
	fmt.Println("running dev-blu")
	fmt.Printf("scriptFile = %+v\n", scriptFile)

	f, err := os.Open(scriptFile)
	if err != nil {
		return
	}

	scriptBytes, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

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
		Script:   string(scriptBytes),
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

		// fmt.Printf("result = %+v\n", result)
		// fmt.Printf("result = %T\n", result)
		// fmt.Println("result", result.String())

		fmt.Println(result.Message)
		// if logM := result.GetLog(); logM != nil {
		// } else if errM := result.GetError(); errM != nil {
		// fmt.Println(errM.Message)
		// fmt.Println("stack:", errM.Stack)
		// } else if completeM := result.GetComplete(); completeM != nil {
		// fmt.Println("done:", completeM.Message)
		// break
		// } else {
		// fmt.Println("result", result.String())
		// }
		if completeM := result.GetComplete(); completeM != nil {
			break
		}
	}

	return
}
