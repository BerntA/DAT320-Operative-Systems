// +build !solution

// Leave an empty line above this comment.
package zlog

import (
	"fmt"
	"time"

	. "github.com/uis-dat320-fall18/printf"
)

type Zaps []ChZap

func NewSimpleZapLogger() ZapLogger {
	zs := make(Zaps, 0)
	return &zs
}

func (zs *Zaps) LogZap(z ChZap) {
	*zs = append(*zs, z)
}

func (zs *Zaps) LogStatus(z StatusChange) {
	// Not implemented for this logger type.
}

func (zs *Zaps) Entries() int {
	return len(*zs)
}

func (zs *Zaps) String() string {
	return fmt.Sprintf("SS: %d", len(*zs))
}

// Viewers() returns the current number of viewers for a channel.
func (zs *Zaps) Viewers(chName string) int {
	defer TimeElapsed(time.Now(), "simple.Viewers")

	viewers := 0
	for _, v := range *zs {
		if v.ToChan == chName {
			viewers++
		}
		if v.FromChan == chName {
			viewers--
		}
	}

	return viewers
}

// Channels() creates a slice of the channels found in the zaps(both to and from).
func (zs *Zaps) Channels() []string {
	defer TimeElapsed(time.Now(), "simple.Channels")

	chanMap := make(map[string]interface{}, 0)
	for _, v := range *zs {
		chanMap[v.ToChan] = nil
		chanMap[v.FromChan] = nil
	}

	chanSlice := make([]string, len(chanMap))
	for k := range chanMap {
		chanSlice = append(chanSlice, k)
	}

	return chanSlice
}

// ChannelsViewers() creates a slice of ChannelViewers, which is defined in zaplogger.go.
// This is the number of viewers for each channel.
func (zs *Zaps) ChannelsViewers() []*ChannelViewers {
	defer TimeElapsed(time.Now(), "simple.ChannelsViewers")

	channelViewers := make([]*ChannelViewers, 0)
	channels := zs.Channels()

	for _, ch := range channels {
		channelViewers = append(channelViewers, &ChannelViewers{Channel: ch, Viewers: zs.Viewers(ch)})
	}

	return channelViewers
}

func (zs *Zaps) AverageChannelDurations() []*ChannelDuration {
	return nil
}

func (zs *Zaps) GetTopMutedChannelDurations() []*ChannelMuteDuration {
	return nil
}
