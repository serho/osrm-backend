package trafficproxyclient

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	maxMsgSize              = 1024 * 1024 * 1024
	defaultTrafficProxyPort = 10086
)

// newGRPCConnection create a new GRPC connection to target traffic proxy.
func newGRPCConnection() (*grpc.ClientConn, error) {

	// make RPC client
	targetServer := flags.target
	if !strings.Contains(targetServer, ":") {
		targetServer = flags.target + ":" + strconv.Itoa(defaultTrafficProxyPort)
	}
	glog.Infoln("dialing traffic proxy " + targetServer)
	conn, err := grpc.Dial(targetServer, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return nil, fmt.Errorf("fail to dial traffic proxy %s, err %v", targetServer, err)
	}
	return conn, nil
}
