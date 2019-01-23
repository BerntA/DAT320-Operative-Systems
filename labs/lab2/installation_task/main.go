// +build !solution

// Leave an empty line above this comment.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/uis-dat320-fall18/assignments/lab2/installation_task/config"
)

var (
	Nameflaged   = flag.String("name", "hello", "The name you will store.")
	Numberflaged = flag.Int("number", 1, "The number you want to store.")
)

func main() {
	flag.Parse()

	cfg1 := config.Configuration{*Numberflaged, *Nameflaged}
	if err := cfg1.Save(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cfg2 := config.Configuration{*Numberflaged, *Nameflaged}
	if err := cfg2.SaveGob(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(cfg1)
	fmt.Println(cfg2)

	// TODO: Load Configuration objects back from disk

	// Loading text file.
	cfg1_loaded, err := config.LoadConfig("config.txt")

	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	} else {
		if cfg1_loaded.Name == "hello" /**Nameflaged*/ && cfg1_loaded.Number == 1 /**Numberflaged*/ {
			fmt.Println("Text parsed correctly.")
		}
	}

	// Loading gob file
	cfg2_loaded, err2 := config.LoadGobConfig("config.gob")

	if err2 != nil {
		fmt.Println(err2)
		//os.Exit(1)
	} else {
		if cfg2_loaded.Name == *Nameflaged && cfg2_loaded.Number == *Numberflaged {
			fmt.Println("Gob parsed correctly.")
		}
	}

	// Printing result.
	fmt.Println(*cfg1_loaded)
	fmt.Println(*cfg2_loaded)
}
