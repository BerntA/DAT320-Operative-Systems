package zlog

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	. "github.com/uis-dat320-fall18/printf"
)

// ZapData This is the data we store for each customer.
type ZapData struct {
	CurrentChannel  *ChZap
	PreviousChannel *ChZap
	MuteState       int
	MuteStart       time.Time
	MutedChannel    string
}

// ChannelData consist of data per TV channel.
// ActiveViewers represents the active viewer count right now
// TotalViews is incremented for every channel swap to this channel
// WatchedFor is the total time that everyone has stayed on this channel for,
// if the value is 60, people have been watching this channel for a total of 60 sec,
// we can then easily compute the avg viewer count for this channel by simply doing TotalViews / WatchedFor,
// everytime someone swaps channel, we increase the previous channel's WatchedFor by the amount of time 'watched'...
type ChannelData struct {
	ActiveViewers int
	TotalViews    int64
	WatchedFor    time.Duration

	MutedViewers    int               // Current muted viewers for this channel.
	TotalMutedViews int64             // Total registered mutes for this channel, incremented during channel swap and mute state changes.
	MuteWatchedFor  time.Duration     // Total mute duration for this channel.
	MutedStatistics []*MutedViewsData // How many muted viewers at a certain time during the day.
}

func (c *ChannelData) UpdateCheckMutes(added bool, t time.Duration) {
	if added {
		c.MutedViewers++
	} else {
		c.MutedViewers--
		c.TotalMutedViews++
		c.MuteWatchedFor += t
	}

	c.MutedStatistics = append(c.MutedStatistics,
		&MutedViewsData{Views: c.MutedViewers, Time: time.Now().Format("15:04:05")})
}

// EffectiveLogger Contains all the necessary data.
type EffectiveLogger struct {
	customers map[string]*ZapData
	channels  map[string]*ChannelData
	lock      *sync.Mutex
}

// NewEffectiveZapLogger Creates and returns a pointer to the effective logger, implementing the ZapLogger interface.
func NewEffectiveZapLogger() ZapLogger {
	zs := &EffectiveLogger{
		customers: make(map[string]*ZapData, 0),
		channels:  make(map[string]*ChannelData, 0),
		lock:      &sync.Mutex{}}
	return zs
}

func (zs *EffectiveLogger) TryRegisterChannel(s string) *ChannelData {
	v, ok := zs.channels[s]
	if ok {
		return v
	}

	data := &ChannelData{MutedStatistics: make([]*MutedViewsData, 0)}
	zs.channels[s] = data
	return data
}

// FindCustomer - Finds an existing customer or creates a new one if not found.
func (zs *EffectiveLogger) FindCustomer(s string) *ZapData {
	zs.lock.Lock()
	defer zs.lock.Unlock()

	v, ok := zs.customers[s]
	if ok {
		return v
	}

	data := &ZapData{}
	zs.customers[s] = data
	return data
}

// LogZap logs a channel zap for the desired customer IP.
func (zs *EffectiveLogger) LogZap(z ChZap) {
	entry := zs.FindCustomer(z.IP.String())
	if entry == nil {
		return
	}

	zs.lock.Lock()
	defer zs.lock.Unlock()

	entry.PreviousChannel = entry.CurrentChannel
	entry.CurrentChannel = &ChZap{IP: z.IP, Time: z.Time, ToChan: z.ToChan, FromChan: z.FromChan}

	if entry.PreviousChannel != nil {
		dur := entry.CurrentChannel.Duration(*entry.PreviousChannel)
		zs.channels[entry.PreviousChannel.ToChan].ActiveViewers--
		zs.channels[entry.PreviousChannel.ToChan].TotalViews++
		zs.channels[entry.PreviousChannel.ToChan].WatchedFor += dur
		//fmt.Printf("Duration since last event for '%v' is '%v'\n", z.IP.String(), dur)

		// We changed channel, check if we are muted, if so also apply mute to this new channel.
		// Store the mute duration, amount of muters for this prev. channel, decrease active muters, etc...
		if (entry.MuteState >= 1) && (entry.MutedChannel != "") {
			zs.channels[entry.MutedChannel].UpdateCheckMutes(false, time.Now().Sub(entry.MuteStart))
		}
	}

	// Guarantees a non-nil accession to the underlying ChannelData struct for the channel name.
	zs.TryRegisterChannel(entry.CurrentChannel.ToChan).ActiveViewers++

	if entry.MuteState >= 1 {
		entry.MutedChannel = entry.CurrentChannel.ToChan
		entry.MuteStart = time.Now()
		zs.channels[entry.CurrentChannel.ToChan].UpdateCheckMutes(true, time.Duration(0))
	}
}

// LogStatus logs a statuschange for the desired customer IP.
func (zs *EffectiveLogger) LogStatus(z StatusChange) {
	entry := zs.FindCustomer(z.IP.String())
	if entry == nil {
		return
	}

	zs.lock.Lock()
	defer zs.lock.Unlock()

	if strings.Compare(z.Type, "Mute_Status") == 0 {
		if entry.MuteState != z.Value {
			if z.Value >= 1 {
				entry.MuteStart = time.Now()
				if entry.CurrentChannel != nil {
					entry.MutedChannel = entry.CurrentChannel.ToChan
					zs.channels[entry.MutedChannel].UpdateCheckMutes(true, time.Duration(0))
				} else {
					entry.MutedChannel = ""
				}
			} else {
				if entry.MutedChannel != "" {
					zs.channels[entry.MutedChannel].UpdateCheckMutes(false, time.Now().Sub(entry.MuteStart))
				}
				entry.MutedChannel = ""
			}
		}
		entry.MuteState = z.Value
	}
}

// Entries returns the amount of customers logged in this logger.
func (zs *EffectiveLogger) Entries() int {
	zs.lock.Lock()
	defer zs.lock.Unlock()
	return len(zs.customers)
}

func (zs *EffectiveLogger) String() string {
	return fmt.Sprintf("SS: %d", zs.Entries())
}

// Viewers returns the current number of viewers for a channel.
func (zs *EffectiveLogger) Viewers(chName string) int {
	zs.lock.Lock()
	defer zs.lock.Unlock()
	//defer TimeElapsed(time.Now(), "efficient.Viewers")

	v, ok := zs.channels[chName]
	if ok == false {
		return 0
	}

	return v.ActiveViewers
}

// Channels creates a slice of the channels found in the logger.
func (zs *EffectiveLogger) Channels() []string {
	zs.lock.Lock()
	defer zs.lock.Unlock()
	//defer TimeElapsed(time.Now(), "efficient.Channels")

	chanSlice := make([]string, len(zs.channels))
	for k := range zs.channels {
		chanSlice = append(chanSlice, k)
	}

	return chanSlice
}

// ChannelsViewers returns a slice with all the channels and the viewer count for each respective channel.
func (zs *EffectiveLogger) ChannelsViewers() []*ChannelViewers {
	defer TimeElapsed(time.Now(), "efficient.ChannelsViewers")

	channelViewers := make([]*ChannelViewers, 0)
	channels := zs.Channels()

	for _, ch := range channels {
		channelViewers = append(channelViewers, &ChannelViewers{Channel: ch, Viewers: zs.Viewers(ch)})
	}

	return channelViewers
}

// AverageChannelDurations returns the most 'switched' channels,
// channels which receive a large duration = many people potentially switched from that channel,
// channels that have 0 duration = no one ever switched from that channel.
func (zs *EffectiveLogger) AverageChannelDurations() []*ChannelDuration {
	zs.lock.Lock()
	defer zs.lock.Unlock()
	defer TimeElapsed(time.Now(), "efficient.AverageChannelDurations")

	avgViewDurPerChannel := make([]*ChannelDuration, 0)
	for k, v := range zs.channels {
		if v.TotalViews <= 0 || v.WatchedFor <= 0 {
			continue
		}

		avg := v.WatchedFor / time.Duration(v.TotalViews)
		if avg <= 0 {
			continue
		}

		avgViewDurPerChannel = append(avgViewDurPerChannel, &ChannelDuration{Channel: k, Duration: avg})
	}

	return avgViewDurPerChannel
}

func (zs *EffectiveLogger) GetTopMutedChannelDurations() []*ChannelMuteDuration {
	zs.lock.Lock()
	defer zs.lock.Unlock()
	defer TimeElapsed(time.Now(), "efficient.GetTopMutedChannelDurations")

	dataDurations := make([]*ChannelMuteDuration, 0)
	for k, v := range zs.channels {
		if v.MuteWatchedFor <= 0 || v.TotalMutedViews <= 0 || len(v.MutedStatistics) <= 0 {
			continue
		}

		avg := v.MuteWatchedFor / time.Duration(v.TotalMutedViews)
		if avg <= 0 {
			continue
		}

		dataDurations = append(dataDurations,
			&ChannelMuteDuration{Channel: k, Duration: avg, Viewers: v.MutedStatistics})
	}

	sort.Slice(dataDurations, func(i, j int) bool {
		return (dataDurations[i].Duration > dataDurations[j].Duration)
	})

	return dataDurations
}
