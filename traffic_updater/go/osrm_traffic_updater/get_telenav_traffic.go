package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Telenav/osrm-backend/traffic_updater/go/gen-go/proxy"
	"github.com/apache/thrift/lib/go/thrift"
)

func getTrafficFlow(ip string, port int, flows []*proxy.Flow, c chan<- bool) ([]*proxy.Flow) {
	var transport thrift.TTransport
	var err error

	startTime := time.Now()

	// make socket
	targetServer := ip + ":" + strconv.Itoa(port)
	fmt.Println("connect traffic proxy " + targetServer)
	transport, err = thrift.NewTSocket(targetServer)
	if err != nil {
		fmt.Println("Error opening socket:", err)
		c <- false
		return nil
	}

	// Buffering
	transport, err = thrift.NewTFramedTransportFactoryMaxLength(thrift.NewTTransportFactory(), 1024*1024*1024).GetTransport(transport)
	if err != nil {
		fmt.Println("Error get transport:", err)
		c <- false
		return nil
	}
	defer transport.Close()
	if err := transport.Open(); err != nil {
		fmt.Println("Error opening transport:", err)
		c <- false
		return nil
	}

	// protocol encoder&decoder
	protocol := thrift.NewTCompactProtocolFactory().GetProtocol(transport)

	// create proxy client
	client := proxy.NewProxyServiceClient(thrift.NewTStandardClient(protocol, protocol))

	// get flows
	fmt.Println("getting flows")
	var defaultCtx = context.Background()
	flows, err = client.GetAllFlows(defaultCtx)
	if err != nil {
		fmt.Println("get flows failed:", err)
		c <- false
		return nil
	}
	fmt.Printf("got flows count: %d\n", len(flows))
	c <- true

	endTime := time.Now()
	fmt.Printf("Processing time for get traffic flow takes %f seconds\n", endTime.Sub(startTime).Seconds())

	return flows
}
