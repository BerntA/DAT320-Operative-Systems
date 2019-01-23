package config

import (
	"encoding/gob"
	"os"
)

func (c Configuration) SaveGob() error {
	cfile, err := os.Create("config.gob")
	if err != nil {
		// couldn't create file
		return err
	}
	// successfully created file; close it at the end of this function
	defer cfile.Close()
	// create a gob encoder for the file
	encoder := gob.NewEncoder(cfile)
	err = encoder.Encode(c)
	return err
}
