package utils

import (
	"log"
	"time"
)

// TimeElapsed measures the time it takes to execute a function.
// Use it as like this with defer:
//     defer TimeElapsed(time.Now(), "FunctionToTime")
//
// For more details see: https://coderwall.com/p/cp5fya
func TimeElapsed(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %v", name, elapsed)
}

const (
	GRPC_SUBSCRIPTION_TYPE_CHANNELS  = 1
	GRPC_SUBSCRIPTION_TYPE_DURATIONS = 2
	GRPC_SUBSCRIPTION_TYPE_SMA       = 3
	GRPC_SUBSCRIPTION_TYPE_MUTED     = 4
)
