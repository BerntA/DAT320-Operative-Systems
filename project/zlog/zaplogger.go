package zlog

import (
	"fmt"
	"time"

	. "github.com/uis-dat320-fall18/printf"
)

type ZapLogger interface {
	LogZap(z ChZap)
	LogStatus(z StatusChange)
	Entries() int
	Viewers(channelName string) int
	Channels() []string
	ChannelsViewers() []*ChannelViewers
	AverageChannelDurations() []*ChannelDuration
	GetTopMutedChannelDurations() []*ChannelMuteDuration
}

type ChannelDuration struct {
	Channel  string
	Duration time.Duration
}

func (cd ChannelDuration) String() string {
	return fmt.Sprintf("%s: %v", cd.Channel, cd.Duration)
}

// MutedViewsData - The muted views at a specific time during the day.
type MutedViewsData struct {
	Views int
	Time  string
}

type ChannelMuteDuration struct {
	Channel  string
	Duration time.Duration
	Viewers  []*MutedViewsData
}

type ChannelViewers struct {
	Channel string
	Viewers int
}

func (cv ChannelViewers) String() string {
	return fmt.Sprintf("%s: %d", cv.Channel, cv.Viewers)
}

type ChanViewersList []*ChannelViewers
type ByViewers struct{ ChanViewersList }

func (t ChanViewersList) Len() int      { return len(t) }
func (t ChanViewersList) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (s ByViewers) Less(i, j int) bool {
	return s.ChanViewersList[i].Viewers < s.ChanViewersList[j].Viewers
}
