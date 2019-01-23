// Just a rough randomizer for testing locally.

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CMD_CHANNEL_SWAP  = 0
	CMD_STATUS_MUTE   = 1
	CMD_STATUS_VOLUME = 2
	CMD_STATUS_HDMI   = 3
	CMD_EVENT_MAX     = 4
)

type ZapEvent struct {
	cmd  int
	info string
}

func (s *ZapEvent) String() string {
	switch s.cmd {
	case CMD_CHANNEL_SWAP:
		return s.info
	case CMD_STATUS_VOLUME:
		return fmt.Sprintf("Volume: %v", rand.Intn(101))
	case CMD_STATUS_MUTE:
		return fmt.Sprintf("Mute_Status: %v", rand.Intn(2))
	case CMD_STATUS_HDMI:
		return fmt.Sprintf("HDMI_Status: %v", rand.Intn(2))
	}

	return "ERROR"
}

var ChannelPool = [...]string{
	"MAX",
	"Viasat 4",
	"FEM",
	"Canal 9",
	"TV2 Film",
	"TV2 Bliss",
	"NRK1",
	"NRK2",
	"NRK3",
	"TV2",
	"TV3",
	"Comedy Central",
	"FOX",
	"nickelodeon",
	"TVNORGE",
	"TV2 Zebra",
	"TV2 Humor",
}

var IPAddressPool = [...]string{
	"1.213.211.140",
	"194.3.94.96",
	"28.11.80.97",
	"230.99.178.86",
	"31.229.130.220",
	"183.125.204.139",
	"24.9.255.243",
	"21.228.105.38",
	"254.139.186.92",
	"83.81.16.149",
	"249.34.119.30",
	"17.185.250.1",
	"144.22.229.97",
	"35.86.40.88",
	"60.142.61.5",
	"81.93.86.242",
	"172.30.45.135",
	"164.159.229.176",
	"117.132.224.218",
	"29.202.132.21",
	"5.120.85.92",
	"38.2.255.117",
	"240.248.65.237",
	"10.239.120.158",
	"83.44.49.201",
	"130.252.153.187",
	"199.85.113.146",
	"242.107.148.200",
	"71.36.116.36",
	"33.129.114.53",
}

func getRandomChannelSwap() string {
	size := len(ChannelPool)
	channel1 := ChannelPool[rand.Intn(size)]
	channel2 := ChannelPool[rand.Intn(size)]
	for strings.Compare(channel1, channel2) == 0 {
		channel2 = ChannelPool[rand.Intn(size)]
	}
	return fmt.Sprintf("%v, %v", channel1, channel2)
}

func getTimeVar(t int) string {
	if t < 10 {
		return fmt.Sprintf("0%v", t)
	}
	return fmt.Sprintf("%v", t)
}

func getRandomMessage() string {
	commandToRun := CMD_CHANNEL_SWAP // 70 % chance to perform channel swap.
	if rand.Intn(100) <= 30 {        // 30 % chance to perform volume, hdmi or mute changes.
		commandToRun = 1 + rand.Intn(CMD_EVENT_MAX-1)
	}

	date := "2018/10/15"
	tim := fmt.Sprintf("%v:%v:%v", getTimeVar(rand.Intn(24)), getTimeVar(rand.Intn(60)), getTimeVar(rand.Intn(60)))
	event := &ZapEvent{cmd: commandToRun}
	if event.cmd == CMD_CHANNEL_SWAP {
		event.info = getRandomChannelSwap()
	}
	return fmt.Sprintf("%v, %v, %v, %v", date, tim, IPAddressPool[rand.Intn(len(IPAddressPool))], event.String())
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", "224.0.1.130:10000")
	if err != nil {
		fmt.Println(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println(err)
	}

	shouldPrint := false
	inn := bufio.NewScanner(os.Stdin)
	fmt.Print("Print output messages: Y/N ? ")
	for inn.Scan() {
		shouldPrint = (strings.Compare(strings.ToUpper(inn.Text()), "Y") == 0)
		break
	}

	conn.SetWriteBuffer(1024)
	for {
		msg := getRandomMessage()
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println(err)
		} else {
			if shouldPrint {
				fmt.Printf("Sending Msg: %v\n", msg)
			}
		}
		time.Sleep(time.Millisecond * 50)
	}
}
