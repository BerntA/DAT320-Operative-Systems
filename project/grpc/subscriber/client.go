package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/uis-dat320-fall18/printf"
	pb "github.com/uis-dat320-fall18/printf/grpc/proto"
	"google.golang.org/grpc"
)

var (
	defaultRefreshRate      = flag.Uint("rate", 1, "Initial refresh rate")
	defaultSubscriptionType = flag.Uint("type", utils.GRPC_SUBSCRIPTION_TYPE_CHANNELS, "Initial subscription type\n0 = Top 10 Channels, 1 = Duration Statistics")
	defaultWindowLength     = flag.Uint("window", 10, "Initial SMA window length")
	defaultChannel          = flag.String("channel", "NRK1", "Initial SMA channel")
	showHelp                = flag.Bool("h", false, "Show this help message and exit")
)

var wantedType uint
var wantedRate uint
var wantedSMAChannel string
var wantedSMAWindow uint

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
		return
	}

	wantedType = *defaultSubscriptionType
	wantedRate = *defaultRefreshRate
	wantedSMAChannel = *defaultChannel
	wantedSMAWindow = *defaultWindowLength
}

// Send to publisher, allow updating the refresh rate.
func sendToPublisher(sub *pb.Subscription_SubscribeClient) {
	(*sub).Send(&pb.SubscribeMessage{Rate: uint32(wantedRate), Type: uint32(wantedType), Window: uint32(wantedSMAWindow), Channel: wantedSMAChannel})
	fmt.Printf("\nTo alter the type, rate, SMA Channel or SMA Window, use these commands:\ntype <arg>\nrate <arg>\nchannel <arg>\nwindow <arg>\n\n")
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		txt := scan.Text()
		cmds := strings.Split(txt, " ")
		if len(cmds) < 2 {
			continue
		}

		val, err := strconv.Atoi(cmds[1])
		if strings.Compare(cmds[0], "channel") != 0 && err != nil {
			fmt.Println(err)
			continue
		}

		switch cmds[0] {
		case "type":
			wantedType = uint(val)
			fmt.Printf("Updated type to %v!\n", wantedType)
		case "rate":
			if val <= 0 {
				fmt.Println("Refresh rate must be greater than zero!")
				continue
			}
			wantedRate = uint(val)
			fmt.Printf("Updated rate to %v sec!\n", wantedRate)
		case "channel":
			wantedSMAChannel = strings.Join(cmds[1:], "")
			fmt.Printf("Updated SMA Channel to %v!\n", wantedSMAChannel)
		case "window":
			if val <= 0 {
				fmt.Println("SMA Window Interval must be greater than zero!")
				continue
			}
			wantedSMAWindow = uint(val)
			fmt.Printf("Updated SMA Window to %v sec!\n", wantedSMAWindow)
		}

		(*sub).Send(&pb.SubscribeMessage{Rate: uint32(wantedRate), Type: uint32(wantedType), Window: uint32(wantedSMAWindow), Channel: wantedSMAChannel})
	}
}

// Receive from publisher, for ex top 10 channels list.
func recvFromPublisher(sub *pb.Subscription_SubscribeClient) {
	for {
		data, err := (*sub).Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println(err)
			break
		}

		if data.Data != nil && len(data.Data) > 0 {
			for i, v := range data.Data {
				fmt.Printf("%v - %v\n", (i + 1), v)
			}
		}

		if data.Sma != nil && len(data.Sma) > 0 {
			avg := float64(0)
			for _, v := range data.Sma {
				avg += float64(v.Views)
			}
			avg /= float64(len(data.Sma))
			fmt.Printf("Average viewers on %v is %v.\n", wantedSMAChannel, math.Round(avg))
		}

		if data.Topmuted != nil && len(data.Topmuted) > 0 {
			for i, v := range data.Topmuted {
				fmt.Printf("%v - %v, avg muted %v, views %v at %v.\n", (i + 1), v.Channel, v.Duration, v.Views, v.Time)
			}
		}

		fmt.Println("")
	}
}

func main() {
	fmt.Println("Started GRPC Client")
	parseFlags()
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
		return
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Kill, os.Interrupt)

	defer conn.Close()

	ctx := context.Background()
	sub, err := pb.NewSubscriptionClient(conn).Subscribe(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	go sendToPublisher(&sub)
	go recvFromPublisher(&sub)

	<-signalChan
}
