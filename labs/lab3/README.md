![UiS](https://www.uis.no/getfile.php/13391907/Biblioteket/Logo%20og%20veiledninger/UiS_liggende_logo_liten.png)

# Lab 3: Network Programming in Go

| Lab 3:		| Network Programming in Go		|
| -------------------- 	| ------------------------------------- |
| Subject: 		| DAT320 Operating Systems 		|
| Deadline:		| Sep 28 2018 16:00			|
| Expected effort:	| 10-15 hours 				|
| Grading: 		| Pass/fail 				|
| Submission: 		| Individually				|

### Table of Contents

1. [Introduction](https://github.com/uis-dat320-fall18/assignments/blob/master/lab3/README.md#introduction)
2. [UDP Echo Server](https://github.com/uis-dat320-fall18/assignments/blob/master/lab3/README.md#udp-echo-server)
3. [Remote Procedure Call](https://github.com/uis-dat320-fall18/assignments/blob/master/lab3/README.md#remote-procedure-call)
4. [Lab Approval](https://github.com/uis-dat320-fall18/assignments/blob/master/lab3/README.md#lab-approval)

## Introduction

This lab should get you started on network programming in Go and introduce some of the
structs and functions in the net package. We recommend browsing over the
[net package documentation](http://golang.org/pkg/net).

**Network questions**
```
1.) Which best describes TCP?
	a. Reliable, connection-oriented
	b. Reliable, connectionless
	c. Unreliable, connection-oriented
	d. Unreliable, connectionless

2.) What is the difference between the types IPAddr and TCPAddr?
	a. IPAddr has variable IP.
	b. TCPAddr has variable IP.
	c. IPAddr has variable Port.
	d. TCPAddr has variable Port.

3.) Which IP address:port is invalid?
	a. 152.94.1.128:65535
	b. 152.94.0.128:12110
	c. 152.94.1.256:12110
	d. 152.94.1.128:12110

4.) Which is not true concerning Gob?
	a. Gob serializes data.
	b. If a field is present in the sender but missing in the receiver, an error will be created.
	c. Gob includes type information. It is self describing.
	d. Gob does not support channels.

5.) Which is a valid RPC declaration in Go? This is from the net/rpc package, not gRPC.
	a. func (t *TypeName) MethodName(args *Args, reply *int) error
	b. func (t *TypeName) MethodName(args *Args) (reply *int, error)
	c. func (t *TypeName) MethodName(args *Args) reply *int
	d. func (t *TypeName) MethodName(args *Args, reply *int)
```

## UDP Echo Server

In this task we will focus on the user datagram protocol (UDP), which provides
unreliable datagram service. You will find the documentation of the
[UDPConn](https://golang.org/pkg/net/#UDPConn) type useful.

In the provided code under `uecho`, we have implemented a simple
`SendCommand()` function that acts as a client, along with a bunch of tests.
You can run these test with `go test -v`, and as described in Lab 1, you can
use the `-run` flag to run only a specific test.

You can also compile your server code into a binary using `go build`. This
will produce a file named `uecho` in the same folder as the `.go` source files.
You can run this binary in two ways:

1. `./uecho -server &` will start the server in the background. Note: *This will
   not work until you have implemented the necessary server parts.*

2. `./uecho` will start the command line client, from which you may interact with
   the server by typing commands into the terminal window.

If you want to extend the capabilities of this runnable client and server,
you can edit the files `echo.go` and `echo_client.go`. But note that the
tests executed by the autograder will use original `SendCommand()` provided
in the original `echo_client.go` file. If you've done something fancy,
and want to show us that's fine, but it won't be considered by the autograder.

#### Echo server specification:


The `SendCommand()` takes the following arguments:

| Argument | Description	|
| -------------------- 	| ------------------------------------- |
| `udpAddr`		| UDP address of the server (`localhost:12110`) 		|
| `cmd`			| Command (as a text string) that the server should interpret and execute |
| `txt`			| Text string on which the server should perform the command provided in `cmd` |

The `SendCommand()` function produces a string composed of the following

```
cmd|:|txt
```

For example:

```
UPPER|:|i want to be upper case
```

From this, the server is expected to produce the following reply:

```
I WANT TO BE UPPER CASE
```

See below for more details about the specific behaviors of the server.

1. For each of the following commands, implement the corresponding functions, so that the returned value corresponds to the expected test outcome. Here you are expected to implement demultiplexer that demultiplexes the input (the command) so that different actions can be taken. A hint is to use the `switch` statement. You will probably also need the `strings.Split()` function.

    | Command	| Action |
    | -------------------- 	| ------------------------------------- |
    | UPPER		| Takes the provided input string `txt` and applies the translates it to upper case using `strings.ToUpper()`. |
    | LOWER		| Same as UPPER, but lower case instead. |
    | CAMEL		| Same as UPPER, but title or camel case instead. |
    | ROT13		| Takes the provided input string `txt` and applies the rot13 translation to it; see lab1 for an example. |
    | SWAP		| Takes the provided input string `txt` and inverts the case. For this command you will find the `strings.Map()` function useful, together with the `unicode.IsUpper()` and `unicode.ToLower()` and a few other similar functions. |

2. The server should reply `Unknown command` if it receives an unknown command
   or fails to interpret a request in any way.

3. Make sure that your server continues to function even if one client's
   connection or datagram packet caused an error.

#### Echo server implementation

You should implement the specification by extending the skeleton code found in
`echo_server.go`:

```go
// +build !solution

// Leave an empty line above this comment.
package main

import (
	"net"
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
	// TODO(student): Implement
	return nil, nil
}

// ServeUDP starts the UDP server's read loop. The server should read from its
// listening socket and handle incoming client requests as according to the
// the specification.
func (u *UDPServer) ServeUDP() {
	// TODO(student): Implement
}

// socketIsClosed is a helper method to check if a listening socket has been
// closed.
func socketIsClosed(err error) bool {
	if strings.Contains(err.Error(), "use of closed network connection") {
		return true
	}
	return false
}
```

## Remote Procedure Call

A popular way to design distributed applications is by means of remote procedure calls. gRPC is
Google's remote procedure calls. You can read about remote procedure calls and gRPC here
[gRPC: Getting Started](https://grpc.io/docs/guides/index.html),
[gRPC: Concepts](https://grpc.io/docs/guides/concepts.html),
[grpc: Go Quickstart](https://grpc.io/docs/quickstart/go.html) and here
[gRPC Basics: Go](https://grpc.io/docs/tutorials/basic/go.html)

It would be useful to read about protocol buffers, Google's way of serializing data. 
[Protocal Buffers: Developer Guide](https://developers.google.com/protocol-buffers/docs/overview)

**Run the following commands to install the necessary libraries.**
```
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

go get -u google.golang.org/grpc
```

There is an example program in grpc called route_guide. Run the server and the client to test your 
installation. Open two terminals. In one run the following commands:
```
cd $GOPATH/src/google.golang.org/grpc/examples/route_guide/

go run server/server.go
```

If the server runs successfully, you should not see any output on this terminal. If you get an error, ask your teaching staff.

In another terminal window, type the following to start a client that will connect to the server that you just started:
```
cd $GOPATH/src/google.golang.org/grpc/examples/route_guide/

go run client/client.go
```

If the client runs successfully, you should see a lot of output similar to this:
```
...
2015/08/10 15:16:10 location:<latitude:405002031 longitude:-748407866 > 
2015/08/10 15:16:10 location:<latitude:409532885 longitude:-742200683 > 
2015/08/10 15:16:10 location:<latitude:416851321 longitude:-742674555 > 
2015/08/10 15:16:10 name:"3387 Richmond Terrace, Staten Island, NY 10303, USA" location:<latitude:406411633 longitude:-741722051 > 
2015/08/10 15:16:10 name:"261 Van Sickle Road, Goshen, NY 10924, USA" location:<latitude:413069058 longitude:-744597778 > 
2015/08/10 15:16:10 location:<latitude:418465462 longitude:-746859398 > 
2015/08/10 15:16:10 location:<latitude:411733222 longitude:-744228360 > 
2015/08/10 15:16:10 name:"3 Hasta Way, Newton, NJ 07860, USA" location:<latitude:410248224 longitude:-747127767 > 
2015/08/10 15:16:10 Traversing 74 points.
2015/08/10 15:16:10 Route summary: point_count:74 distance:720194535 
2015/08/10 15:16:10 Got message First message at point(0, 1)
2015/08/10 15:16:10 Got message Fourth message at point(0, 1)
2015/08/10 15:16:10 Got message First message at point(0, 1)
...
```

You may want to study the code for the `route_guide` example to try to understand what is going on.

# Building your own RPC-based key-value storage

RPC key-value store: The repository contains code for a very simple key-value storage
service, where the keys and values are strings. It offers the method Insert, which inserts a
key-value pair. It returns a bool indicating success or failure.

**Please look at, but do not change, the kv.pb.go file. This file contains important APIs and message
definitions needed for these exercises.**


```
type InsertRequest struct {
	Key   string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
}

type InsertResponse struct {
	Success bool `protobuf:"varint,1,opt,name=success" json:"success,omitempty"`
}

type LookupRequest struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
}

type LookupResponse struct {
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

type KeysRequest struct {
}

type KeysResponse struct {
	Keys []string `protobuf:"bytes,1,rep,name=keys" json:"keys,omitempty"`
}

// Client API for KeyValueService service

type KeyValueServiceClient interface {
	Insert(ctx context.Context, in *InsertRequest, opts ...grpc.CallOption) (*InsertResponse, error)
	Lookup(ctx context.Context, in *LookupRequest, opts ...grpc.CallOption) (*LookupResponse, error)
	Keys(ctx context.Context, in *KeysRequest, opts ...grpc.CallOption) (*KeysResponse, error)
}

// Server API for KeyValueService service

type KeyValueServiceServer interface {
	Insert(context.Context, *InsertRequest) (*InsertResponse, error)
	Lookup(context.Context, *LookupRequest) (*LookupResponse, error)
	Keys(context.Context, *KeysRequest) (*KeysResponse, error)
}
```

**Tasks:**

1. Create a client and a connection to the server.

2. In the client, call the Insert() gRPC for a number of key/value pairs.

3. In the server, implement the Lookup() gRPC, which should return the value of the requested key.

4. In the server, implement the Keys() gRPC, which should return a slice of the keys 
(not the values) of the map back to the client.

5. In the client, call the Lookup() gRPC for each of the key/value pairs inserted and verify 
the result returned from the Lookup() gRPC matches the value inserted for the corresponding key.

6. In the client, call the Keys() gRPC and verify that the number of keys returned matches 
the expected number.

7. Several clients connecting to the server may read and write concurrently from the shared
key-value map. This will eventually cause inconsistencies in the map, unless some form of
protection is instituted. Implement locking at the appropriate locations
in the code. See [pkg/sync](http://golang.org/pkg/sync/).

Extras:

  - Explain why the clients can access the map at the server concurrently.
  - If you run your server without protection on the map, are you able to provoke inconsistencies in the map.

Troubleshooting if you get compile errors related to `kv.pb.go`, it may help to recompile the proto file:
```
cd lab3/grpc/proto
protoc --go_out=plugins=grpc:. kv.proto
```

## Lab Approval

To have your lab assignment approved, you must come to the lab during lab hours
and present your solution. This lets you present the thought process behind
your solution, and gives us more information for grading purposes. When you are
ready to show your solution, reach out to a member of the teaching staff.  It
is expected that you can explain your code and show how it works. You may show
your solution on a lab workstation or your own computer. The results from
Autograder will also be taken into consideration when approving a lab. At least
60% of the Autograder tests should pass for the lab to be approved. A lab needs
to be approved before Autograder will provide feedback on the next lab
assignment.

Also see the [Grading and Collaboration
Policy](https://github.com/uis-dat320-fall18/course-info/blob/master/policy.md)
document for additional information.
