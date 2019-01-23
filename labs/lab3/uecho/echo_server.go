// +build !solution

// Leave an empty line above this comment.
package main

import (
	"net"
	"strconv"
	"strings"
)

// UDPServer implements the UDP server specification found at
// https://github.com/uis-dat320-fall18/assignments/blob/master/lab3/README.md#echo-server-specification
type UDPServer struct {
	conn *net.UDPConn
	// TODO(student): Add fields if needed
}

// NewUDPServer returns a new UDPServer listening on addr. It should return an
// error if there was any problem resolving or listening on the provided addr.
func NewUDPServer(addr string) (*UDPServer, error) {
	socket := strings.Split(addr, ":")
	port, _ := strconv.Atoi(socket[1])
	address := net.UDPAddr{
		IP:   net.ParseIP(socket[0]),
		Port: port,
	}
	connection, err := net.ListenUDP("udp", &address)
	server := UDPServer{conn: connection}
	return &server, err
}

// ServeUDP starts the UDP server's read loop. The server should read from its
// listening socket and handle incoming client requests as according to the
// the specification.
func (u *UDPServer) ServeUDP() {
	buffer := make([]byte, 1024)
	failMsg := []byte("Unknown command")
	for {
		n, addr, err := u.conn.ReadFromUDP(buffer)
		if err != nil {
			u.conn.WriteToUDP(failMsg, addr)
			continue
		}

		data := strings.Split(string(buffer[0:n]), "|:|")
		if len(data) != 2 {
			u.conn.WriteToUDP(failMsg, addr)
			continue
		}

		switch data[0] {
		case "UPPER":
			u.conn.WriteToUDP([]byte(strings.ToUpper(data[1])), addr)
		case "LOWER":
			u.conn.WriteToUDP([]byte(strings.ToLower(data[1])), addr)
		case "CAMEL":
			u.conn.WriteToUDP([]byte(strings.Title(strings.ToLower(data[1]))), addr)
		case "ROT13":
			u.conn.WriteToUDP([]byte(Rot13(data[1])), addr)
		case "SWAP":
			u.conn.WriteToUDP([]byte(SwapCase(data[1])), addr)
		default:
			u.conn.WriteToUDP(failMsg, addr)
		}

		//fmt.Printf("Data read %v - %v\n", n, string(buffer[0:n]))
	}
}

// socketIsClosed is a helper method to check if a listening socket has been
// closed.
func socketIsClosed(err error) bool {
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}
	return false
}

// Rot13 applies the rot13 algorithm on any string, returns the new string.
func Rot13(s string) string {
	bytes := []byte(s)
	for i, v := range bytes {
		bytes[i] = func(b byte) byte {
			if !(b >= 65 && b <= 90) && !(b >= 97 && b <= 122) {
				return b
			}

			newPos := (b - 13)
			if b >= 65 && b <= 90 {
				if newPos < 65 {
					newPos = 90 - (64 - newPos)
				}
			} else {
				if newPos < 97 {
					newPos = 122 - (96 - newPos)
				}
			}

			return newPos
		}(v)
	}

	return string(bytes[:])
}

// SwapCase will return a new string with each character's case swapped.
func SwapCase(s string) string {
	bytes := []byte(s)
	diff := byte('a' - 'A')
	for i, v := range bytes {
		if v >= 'A' && v < 'a' {
			bytes[i] += diff
		} else if v >= 'a' && v < '{' {
			bytes[i] -= diff
		}
	}

	return string(bytes[:])
}
