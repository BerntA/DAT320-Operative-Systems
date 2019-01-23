// +build !solution

// Leave an empty line above this comment.
package config

import (
	"encoding/gob"
	"os"
)

func LoadGobConfig(file string) (conf *Configuration, err error) {
	// Create a dummy Configuration object into which to decode it
	conf = &Configuration{}
	err = nil
	// TODO: Open file using os.Open()
	// TODO: Decode using a gob decoder.

	// Opening file.
	cfile, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	// Decode using gob.
	decoder := gob.NewDecoder(cfile)
	err = decoder.Decode(conf)

	return conf, err
}
