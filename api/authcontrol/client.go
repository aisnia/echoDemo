package authcontrol

import (
	"context"
	"google.golang.org/grpc"
	pb "learn_together/api/authcontrol/proto"
	"learn_together/service/etcd"
)

type Client struct {
	remoteClient pb.CasbinClient
}

func NewClient(ctx context.Context, address, serviceName string, opts ...grpc.DialOption) (*Client, error) {
	// Set up a connection to the server.
	conn, err := etcd.NewClient(ctx, address, serviceName, opts...)
	if err != nil {
		return nil, err
	}
	c := pb.NewCasbinClient(conn)
	return &Client{
		remoteClient: c,
	}, nil
}
