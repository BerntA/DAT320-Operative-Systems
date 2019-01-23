// +build !solution

// Leave an empty line above this comment.
//
// Zap Collection Server
package main

import (
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	. "github.com/uis-dat320-fall18/printf"
	"github.com/uis-dat320-fall18/printf/zlog"
)

func GetTopTenChannelViews() []string {
	tvChannelsAndViews := ztore.ChannelsViewers()
	sort.Slice(tvChannelsAndViews, func(i, j int) bool {
		return (tvChannelsAndViews[i].Viewers > tvChannelsAndViews[j].Viewers)
	})

	size := len(tvChannelsAndViews)
	channels := make([]string, 0)
	for i := 0; i < 10; i++ {
		if i >= size {
			break
		}
		channels = append(channels, tvChannelsAndViews[i].String())
	}

	return channels
}

func GetChannelDurations() []string {
	durations := ztore.AverageChannelDurations()
	if durations == nil || len(durations) <= 0 {
		return nil
	}

	sort.Slice(durations, func(i, j int) bool {
		return (durations[i].Duration > durations[j].Duration)
	})

	channels := make([]string, 0)
	for _, v := range durations {
		channels = append(channels, v.String())
	}

	return channels
}

func displayViewersForChannel(c string) {
	for {
		fmt.Printf("%v has %v viewers.\n", c, ztore.Viewers(c))
		time.Sleep(1 * time.Second)
	}
}

func displayTopTenChannels() {
	for {
		fmt.Println("Top 10 TV Channels:")

		list := GetTopTenChannelViews()
		if list != nil {
			for i, v := range list {
				fmt.Printf("%v - %v\n", (i + 1), v)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

// REMARK: This function should return (i.e. it should not block)
func runLab() {
	switch *labnum {
	case "a", "c1", "c2", "d", "e":
		ztore = zlog.NewSimpleZapLogger()
	case "f":
		ztore = zlog.NewEffectiveZapLogger()
	}

	if *labnum == "a" {
		go processZapEvents(true, nil)
	} else {
		go processZapEvents(false, &ztore)
	}

	switch *labnum {
	case "c1":
		go displayViewersForChannel("NRK1")
	case "c2":
		go displayViewersForChannel("TV2 Norge")
	case "d":
		// - The viewers are calculated wrong,
		// you store multiple entries for the same customers, causing havoc when calculating the view count.
		// if one customer swapped to TV2, and later swapped to TV Norge, TV2 would always add a -1 view, later a new
		// event could be registered for the same customer where he switches back to TV2, from TV Norge, now it will add a +1 again,
		// BUT the view count will be 0! for both channels. It should be 0 for TV Norge only... And 1 for TV2.
	case "e", "f":
		go displayTopTenChannels()
	}
}

// REMARK: This function should return (i.e. it should not block)
func startServer() {
	log.Println("Starting ZapServer...")
	udpAddress, err := net.ResolveUDPAddr("udp", *maddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	listener, err := net.ListenMulticastUDP("udp", nil, udpAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	server = listener
	server.SetReadBuffer(1024)
	fmt.Printf("Listening to: %v\n", udpAddress)
}

func processZapEvents(shouldPrint bool, logger *zlog.ZapLogger) {
	buffer := make([]byte, 1024)
	for server != nil {
		n, src, err := server.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			break
		}

		event := string(buffer[0:n])
		zapChan, zapStatus, err := NewSTBEvent(event)

		if shouldPrint {
			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Printf("Read Data (%v bytes): %v : From %v\n", n, event, src)
		}

		if logger != nil {
			if zapChan != nil {
				(*logger).LogZap(*zapChan)
				if *grpcPublisher {
					go OnLoggedZapChanForSMA(zapChan.Time.Unix(), zapChan.ToChan, zapChan.FromChan)
				}
			}

			if zapStatus != nil {
				(*logger).LogStatus(*zapStatus)
			}
		}
	}
	fmt.Println("Finished processing ZapEvents...")
}
