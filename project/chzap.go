// +build !solution

// Leave an empty line above this comment.
package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const timeFormat = "2006/01/02, 15:04:05"
const dateFormat = "2006/01/02"
const timeOnly = "15:04:05"
const timeLen = len(timeFormat)

const (
	TYPE_STATUS       = 0
	TYPE_CHANNEL_SWAP = 1
)

type StatusChange struct {
	Time  time.Time
	IP    net.IP
	Type  string
	Value int
}

type ChZap struct {
	Time     time.Time
	IP       net.IP
	ToChan   string
	FromChan string
}

func NewSTBEvent(event string) (*ChZap, *StatusChange, error) {
	data := strings.Split(event, ", ")
	size := len(data)
	if size < 4 {
		if size <= 2 {
			return nil, nil, fmt.Errorf(fmt.Sprintf("NewSTBEvent: too short event string: %v", event))
		}
		return nil, nil, fmt.Errorf(fmt.Sprintf("NewSTBEvent: event with too few fields: %v", event))
	}

	date, err := time.Parse(timeFormat, fmt.Sprintf("%v, %v", data[0], data[1]))
	if err != nil {
		return nil, nil, fmt.Errorf("NewSTBEvent: failed to parse timestamp")
	}

	ipAddress := net.ParseIP(data[2])

	// Check if this is a status event
	if strings.Contains(data[3], ":") {
		info := strings.Split(data[3], ": ")
		val, err := strconv.Atoi(info[1])
		if err != nil {
			return nil, nil, fmt.Errorf("NewSTBEvent: Unable to parse value for status type")
		}

		return nil, &StatusChange{Time: date, IP: ipAddress, Type: info[0], Value: val}, nil
	}

	// Check if ChZap will be valid
	if size < 5 {
		return nil, nil, fmt.Errorf("NewSTBEvent: Unable to parse ChZap, too few arguments")
	}

	// Assuming event is a channel swap
	return &ChZap{Time: date, IP: ipAddress, ToChan: data[3], FromChan: data[4]}, nil, nil
}

func (zap ChZap) String() string {
	return fmt.Sprintf("%v, %v, %v, %v", zap.Time.String(), string(zap.IP[:]), zap.ToChan, zap.FromChan)
}

func (schg StatusChange) String() string {
	return fmt.Sprintf("%v, %v, %v: %v", schg.Time.String(), string(schg.IP[:]), schg.Type, schg.Value)
}

// The duration between receiving (this) zap event and the provided event
func (zap ChZap) Duration(provided ChZap) time.Duration {
	time1, _ := time.Parse(timeOnly, zap.Time.Format(timeOnly))
	time2, _ := time.Parse(timeOnly, provided.Time.Format(timeOnly))
	return time1.Sub(time2)
}

func (zap ChZap) Date() string {
	return zap.Time.Format(dateFormat)
}
