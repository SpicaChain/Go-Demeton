package rpc

import (
	"github.com/cyber-demeton/go-demeton/util/logging"
	"google.golang.org/grpc"
)

// Dial returns a client connection.
func Dial(target string) (*grpc.ClientConn, error) {
	// TODO: support secure connection.
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		logging.VLog().Debug("rpc.Dial() failed: ", err)
	}
	return conn, err
}
