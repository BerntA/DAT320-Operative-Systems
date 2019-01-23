// +build !solution

// Leave an empty line above this comment.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/uis-dat320-fall18/assignments/lab3/grpc/proto"
)

const (
	address = "localhost:12111"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Unable to connect!")
		return
	}

	defer conn.Close()

	kv := pb.NewKeyValueServiceClient(conn)

	ctx := context.Background()
	kv.Insert(ctx, &pb.InsertRequest{Key: "1", Value: "one"})
	kv.Insert(ctx, &pb.InsertRequest{Key: "2", Value: "two"})
	kv.Insert(ctx, &pb.InsertRequest{Key: "Game", Value: "Something"})
	kv.Insert(ctx, &pb.InsertRequest{Key: "Hello", Value: "World"})

	fmt.Println("Started connection successfully!\nCommands:\nL <Key> - Fetch value for some key.\nK - Fetch available keys.\n")

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		cmd := strings.Split(scan.Text(), " ")
		switch cmd[0] {

		case "L":
			if len(cmd) != 2 {
				fmt.Println("Invalid args!")
				continue
			}

			rs, err := kv.Lookup(ctx, &pb.LookupRequest{Key: cmd[1]})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if len(rs.Value) <= 0 {
				fmt.Println("No response!")
				continue
			}

			fmt.Println(rs.Value)

		case "K":
			rs, err := kv.Keys(ctx, &pb.KeysRequest{})
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			if rs.Keys != nil && len(rs.Keys) > 0 {
				fmt.Println(rs.Keys)
			} else {
				fmt.Println("No response!")
			}

		default:
			fmt.Println("Unknown...")

		}
	}
}
