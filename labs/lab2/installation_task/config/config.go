package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Configuration struct {
	Number int
	Name   string
}

// Save configuration c as text file 'config.txt' in the current directory.
func (c Configuration) Save() error {
	cfile, err := os.Create("config.txt")
	if err != nil {
		// couldn't create file
		return err
	}
	// successfully created file; close it at the end of this function
	defer cfile.Close()
	_, err = fmt.Fprintf(cfile, "Number=%d\nName=%s", c.Number, c.Name)
	return err
}

// LoadConfig loads a text-based configuration file, and returns the
// corresponding Configuration object.
func LoadConfig(file string) (*Configuration, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	// create empty Configuration object
	conf := &Configuration{}
	for _, line := range lines {
		err = conf.parse(line)
		if err != nil {
			return nil, err
		}
	}
	return conf, nil
}
