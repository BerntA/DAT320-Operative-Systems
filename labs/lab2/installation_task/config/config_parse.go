// +build !solution

// Leave an empty line above this comment.
package config

import (
	"errors"
	"strconv"
	"strings"
)

func (c *Configuration) parse(line string) (err error) {
	// TODO: find keys=value in line
	//       and store value to the correct part of the c object
	// TIPS: strconv.Atoi()
	a := strings.Split(line, "=")

	if len(a) != 2 {
		return errors.New("Couldn't find a key value pair.")
	}

	if strings.Contains(a[0], "Number") {
		num, err := strconv.Atoi(a[1])

		if err != nil {
			return err
		}

		c.Number = num
	} else if strings.Contains(a[0], "Name") {
		c.Name = a[1]
	} else {
		return errors.New("Didn't contain any correct keys.")
	}

	return nil
}
