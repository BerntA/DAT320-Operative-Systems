// Create a pool of subs @ GRPC...

package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 15; i++ {
		wg.Add(1)
		go func(num int) {
			dur := 20 + rand.Intn(20)
			fmt.Printf("%v Running for %v sec!\n", num, dur)
			cmd := exec.Command("./subscriber/subscriber",
				fmt.Sprintf("-type %v", 1+rand.Intn(5)),
				fmt.Sprintf("-rate %v", 1+rand.Intn(4)))
			cmd.Start()

			time.Sleep(time.Duration(dur) * time.Second)
			cmd.Process.Kill()
			wg.Done()
			fmt.Printf("Finished sub. %v\n", num)
		}(i + 1)
	}
	wg.Wait()
	fmt.Println("Simulators done...")
}
