// +build !solution

package lab2

import (
	"regexp"
	"math/rand"
	"strings"
)

/*
Task 4: Scrambling text

In this task the objective is to transform a given text such that all letters of each word
are randomly shuffled except for the first and last letter.

For example, given the word "scramble", the result could be "srmacble" or "sbcamrle",
or any other permutation as long as the first and last letters stay the same

An entire sentence scrambled like this should still be readable:
"it deosn't mttaer in waht oredr the ltteers in a wrod are,
the olny iprmoetnt tihng is taht the frist and lsat ltteer be at the rghit pclae"
See https://www.mrc-cbu.cam.ac.uk/people/matt.davis/cmabridge/ for more
information and examples.

Implementation:
The task is to implement the scramble function, which takes a text in the form of a string and a seed.
A seed is given so the output from your solution should match the test cases if it is correct.
The seed should be applied at the start of the function.
Remember that the implementation should keep any punctuation and spacing intact, and all numbers should be untouched.

Shuffling the letters and applying the seed can be done using the math/rand package (https://golang.org/pkg/math/rand/).
Use the Shuffle function to ensure you reach the same values as given in the tests (scramle_test.go).

A function for properly tokenizing text is given, since regular expressions are a bit out of the scope of this course.
It will return a slice containing all tokens. In this case, a token may be a word, a space, any punctuation or a number(can be multiple digits).

*/

func splitText(text string) []string {
	re := regexp.MustCompile("[A-Za-z0-9']+|[':;?().,!\\ ]")
	return re.FindAllString(text, -1)
}

func isValid(c byte) bool {
	if !(c >= 65 && c <= 90) && !(c >= 97 && c <= 122) {
		return false
	}
	
	return true
}

func scramble(text string, seed int64) string {
	rnd := rand.New(rand.NewSource(seed))
	tmp := splitText(text)
	var ret string

	first, last := -1, -1

	for _, v := range tmp {
		if len(v) <= 3 {
			ret += v
			continue
		}

		for i := 0; i < len(v); i++ {
			if !isValid(v[i]) {
				continue
			}

			first = i
			last = i
			break
		}

		if first != -1 {
			for i := (len(v) - 1); i > first; i-- {
				if !isValid(v[i]) {
					continue
				}

				last = i			
				break
			}
		}

		if first == -1 && last == -1 {
			ret += v
			continue
		}

		ret += v[0:(first+1)]

		scrambleVal := v[(first + 1) : last]
		wordsToScramble := strings.Split(scrambleVal, "")
		rnd.Shuffle(len(wordsToScramble), func (i, j int) {
			if !isValid(scrambleVal[i]) || !isValid(scrambleVal[j]) {
				return
			}

			wordsToScramble[i], wordsToScramble[j] = wordsToScramble[j], wordsToScramble[i]
		})

		for _, sv := range wordsToScramble {
			ret += sv
		}

		ret += v[last:len(v)]
	}

	return ret
}
