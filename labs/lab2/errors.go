// +build !solution

package lab2

import "fmt"

// Errors is an error returned by a multiwriter when there are errors with one
// or more of the writes. Errors will be in a one-to-one correspondence with
// the input elements; successful elements will have a nil entry.
//
// Should not be constructed if all entires are nil; in this case one should
// instead return just nil instead of a MultiError.
type Errors []error

/*
Task 5: Errors needed for multiwriter

You may find this blog post useful:
http://blog.golang.org/error-handling-and-go

Similar to a the Stringer interface, the error interface also defines a
method that returns a string.

type error interface {
    Error() string
}

Thus also the error type can describe itself as a string. The fmt package (and
many others) use this Error() method to print errors.

Implement the Error() method for the Errors type defined above.

The following conditions should be covered:

1. When there are no errors in the slice, it should return:

"(0 errors)"

2. When there is one error in the slice, it should return:

The error string return by the corresponding Error() method.

3. When there are two errors in the slice, it should return:

The first error + " (and 1 other error)"

4. When there are X>1 errors in the slice, it should return:

The first error + " (and X other errors)"
*/
func (m Errors) Error() string {
	if m == nil {
		return "(0 errors)"
	}

	errorCount := 0
	index := -1
	for i := 0; i < len(m); i++ {
		if (m[i] != nil) {
			errorCount++
			index = i
		}
	}

	if index >= 0 && m[index] != nil {
		if errorCount > 2 {
			return fmt.Sprintf("%v (and %v other errors)", m[index].Error(), (errorCount - 1))
		} else if errorCount == 2 {
			return fmt.Sprintf("%v (and 1 other error)", m[index].Error())
		}
		return m[index].Error()
	}

	return "(0 errors)"
}
