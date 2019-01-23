// Zap Collection Server
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/uis-dat320-fall18/printf/zlog"
)

var (
	maddr         = flag.String("mcast", "224.0.1.130:10000", "multicast ip:port")
	labnum        = flag.String("lab", "c2", "which lab exercise to run")
	grpcPublisher = flag.Bool("grpc", false, "Enable GRPC Pub/Sub Server")
	showHelp      = flag.Bool("h", false, "show this help message and exit")
	memprofile    = flag.String("memprofile", "", "write memory profile to this file")
)

var ztore zlog.ZapLogger
var server *net.UDPConn

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func parseFlags() {
	flag.Usage = Usage
	flag.Parse()
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	parseFlags()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Kill, os.Interrupt)
	startServer()
	if server != nil {
		defer server.Close()
	}
	runLab()

	if *grpcPublisher {
		go StartServerGRPC()
	}

	// Here we wait for CTRL-C or some other kill signal
	s := <-signalChan
	fmt.Println("Server stopping on", s, "signal")
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		fmt.Println("Saved memory profile")
		fmt.Println("Analyze with: go tool pprof $GOPATH/bin/zapserver", *memprofile)
	}
}
