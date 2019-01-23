package main

import (
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/uis-dat320-fall18/printf"
	pb "github.com/uis-dat320-fall18/printf/grpc/proto"
	"google.golang.org/grpc"
)

type PublisherSubscriberServer struct {
	channels  []string
	durations []string
	dataLock  *sync.RWMutex
}

type PublisherSubscriberClient struct {
	typ        uint32
	rate       time.Duration
	smaWindow  uint32
	smaChannel string
	term       chan interface{}
}

func (c *PublisherSubscriberClient) Run(s *PublisherSubscriberServer, stream *pb.Subscription_SubscribeServer) {
	shouldTerm := false
	var timeUntilNextSMA int64
	for {
		if c.rate < 1 || c.typ <= 0 {
			time.Sleep(time.Second * 1)
			continue
		}

		select {
		case <-c.term:
			shouldTerm = true
		default:
		}

		if shouldTerm {
			break
		}

		switch c.typ {
		case utils.GRPC_SUBSCRIPTION_TYPE_CHANNELS:
			(*stream).Send(&pb.NotificationMessage{Data: s.GetTopTenChannels()})
		case utils.GRPC_SUBSCRIPTION_TYPE_DURATIONS:
			(*stream).Send(&pb.NotificationMessage{Data: s.GetChannelDurations()})
		case utils.GRPC_SUBSCRIPTION_TYPE_SMA:
			if timeUntilNextSMA <= 0 {
				timeUntilNextSMA = (time.Now().Unix() + int64(c.smaWindow))
				go c.StartSMAMeasurement(stream)
			} else if time.Now().Unix() > timeUntilNextSMA {
				timeUntilNextSMA = 0
			}
		case utils.GRPC_SUBSCRIPTION_TYPE_MUTED:
			(*stream).Send(&pb.NotificationMessage{Topmuted: s.GetTopTenMutedChannelStatistics()})
		}

		time.Sleep(time.Second * c.rate)
	}
}

func (c *PublisherSubscriberClient) StartSMAMeasurement(stream *pb.Subscription_SubscribeServer) {
	entry := RegisterSMARequest(c.smaChannel, c.smaWindow)
	time.Sleep(time.Second * time.Duration(c.smaWindow))
	(*stream).Send(&pb.NotificationMessage{Sma: entry.GetMeasurements()})
	DeregisterSMARequest(entry)
	entry = nil
}

func (c *PublisherSubscriberClient) Destroy() {
	c.term <- nil
}

func (s *PublisherSubscriberServer) Subscribe(stream pb.Subscription_SubscribeServer) error {
	client := &PublisherSubscriberClient{
		typ:       uint32(0),
		rate:      time.Duration(0),
		smaWindow: uint32(5),
		term:      make(chan interface{}),
	}

	go client.Run(s, &stream)
	defer client.Destroy()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		client.smaChannel = in.GetChannel()
		client.smaWindow = in.GetWindow()
		client.rate = time.Duration(in.GetRate())
		client.typ = in.GetType()
	}

	return nil
}

func (s *PublisherSubscriberServer) GetTopTenChannels() []string {
	s.dataLock.RLock()
	defer s.dataLock.RUnlock()
	return s.channels
}

func (s *PublisherSubscriberServer) GetChannelDurations() []string {
	s.dataLock.RLock()
	defer s.dataLock.RUnlock()
	return s.durations
}

func (s *PublisherSubscriberServer) GetTopTenMutedChannelStatistics() []*pb.TopMutedData {
	data := ztore.GetTopMutedChannelDurations()
	if data == nil || len(data) <= 0 {
		return nil
	}

	defer utils.TimeElapsed(time.Now(), "publisher.GetTopTenMutedChannelStatistics")
	mutedList := make([]*pb.TopMutedData, 0)
	size := len(data)

	for i := 0; i < 10; i++ {
		if i >= size {
			break
		}

		sort.Slice(data[i].Viewers, func(x, y int) bool {
			return (data[i].Viewers[x].Views > data[i].Viewers[y].Views)
		})

		mutedList = append(mutedList,
			&pb.TopMutedData{
				Channel:  data[i].Channel,
				Duration: data[i].Duration.String(),
				Time:     data[i].Viewers[0].Time,
				Views:    uint32(data[i].Viewers[0].Views),
			})
	}

	return mutedList
}

func (s *PublisherSubscriberServer) CacheLoggerDataForGRPC() {
	for {
		topViews := GetTopTenChannelViews()
		topDurations := GetChannelDurations()

		s.dataLock.Lock()
		if topViews != nil {
			s.channels = topViews
		}

		if topDurations != nil {
			s.durations = topDurations
		}
		s.dataLock.Unlock()

		time.Sleep(time.Second * 1)
	}
}

// StartServerGRPC Starts the GRPC service, the server will cache the durations list and top 10 list every sec.
func StartServerGRPC() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Started GRPC Server\n")
	grpcServer := grpc.NewServer()
	pubsubServer := &PublisherSubscriberServer{channels: make([]string, 0), durations: make([]string, 0), dataLock: &sync.RWMutex{}}
	go pubsubServer.CacheLoggerDataForGRPC()
	pb.RegisterSubscriptionServer(grpcServer, pubsubServer)
	grpcServer.Serve(listener)
}
