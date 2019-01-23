package main

import (
	"strings"
	"sync"
	"time"

	pb "github.com/uis-dat320-fall18/printf/grpc/proto"
)

// SMAEntry - SimpleMovingAverage
type SMAEntry struct {
	Start           time.Time
	Interval        float64
	Channel         string               // Channel to monitor
	Viewers         []*pb.SMAMeasurement // Channel views over time
	LastViewerCount uint32               // Last count registered
	mux             *sync.Mutex
}

var movingAverageFuncs = make(map[*SMAEntry]interface{}, 0)
var movingAverageLocker = &sync.Mutex{}

// CheckAddViewers - Check if we can add a viewer measurement for our channel.
func (s *SMAEntry) CheckAddViewers(t int64, to string, from string) {
	// Check if we're done
	dur := time.Now().Sub(s.Start)
	if dur.Seconds() > s.Interval {
		return
	}

	if strings.Compare(s.Channel, to) != 0 && strings.Compare(s.Channel, from) != 0 {
		return
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	views := uint32(ztore.Viewers(s.Channel))
	if views == s.LastViewerCount {
		return
	}

	s.Viewers = append(s.Viewers, &pb.SMAMeasurement{Views: views, Time: t})
	s.LastViewerCount = views
}

// GetMeasurements - Returns the current measurements for this SMA Entry.
func (s *SMAEntry) GetMeasurements() []*pb.SMAMeasurement {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.Viewers
}

// RegisterSMARequest - Register SMA entry.
func RegisterSMARequest(channel string, interval uint32) *SMAEntry {
	movingAverageLocker.Lock()
	defer movingAverageLocker.Unlock()

	v := &SMAEntry{
		Start:           time.Now(),
		Interval:        float64(interval),
		Channel:         channel,
		Viewers:         make([]*pb.SMAMeasurement, 0),
		LastViewerCount: 0,
		mux:             &sync.Mutex{},
	}

	movingAverageFuncs[v] = nil
	return v
}

// DeregisterSMARequest - Removes an SMA entry.
func DeregisterSMARequest(v *SMAEntry) bool {
	movingAverageLocker.Lock()
	defer movingAverageLocker.Unlock()

	_, ok := movingAverageFuncs[v]
	if ok {
		delete(movingAverageFuncs, v)
		return true
	}
	return false
}

// OnLoggedZapChanForSMA - Whenever we receive ANY zapchan event, update our 'listeners', our listeners will check if they care about this channel event,
// if so they will add a new entry with the current viewer count for the desired channel.
func OnLoggedZapChanForSMA(t int64, to string, from string) {
	movingAverageLocker.Lock()
	defer movingAverageLocker.Unlock()

	if movingAverageFuncs == nil || len(movingAverageFuncs) <= 0 {
		return
	}

	for k := range movingAverageFuncs {
		k.CheckAddViewers(t, to, from)
	}
}
