package floodchrome

import (
	"context"

	pb "github.com/flood-io/cli/proto"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	pb.TestClient
}

func NewClient(serverAddr string) (client *Client, err error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return
	}

	client = &Client{
		conn,
		pb.NewTestClient(conn),
	}

	return
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) WsEndpoint() (endpoint string, err error) {
	req := &pb.DevtoolsRequest{}

	resp, err := c.Devtools(context.Background(), req)
	if err != nil {
		return
	}

	endpoint = resp.WsEndpoint
	return
}
