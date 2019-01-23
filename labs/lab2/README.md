![UiS](https://www.uis.no/getfile.php/13391907/Biblioteket/Logo%20og%20veiledninger/UiS_liggende_logo_liten.png)

# Lab 2: Introduction to Go Programming

| Lab 2:		|  Introduction to Go Programming	|
| -------------------- 	| ------------------------------------- |
| Subject: 		| DAT320 Operating Systems 		|
| Deadline:		| Sep 14 2017 16:00			|
| Expected effort:	| 10-15 hours 				|
| Grading: 		| Pass/fail 				|
| Submission: 		| Individually				|

### Table of Contents

1. [Introduction](https://github.com/uis-dat320-fall18/assignments/blob/master/lab2/README.md#introduction)
2. [The Go Language](https://github.com/uis-dat320-fall18/assignments/blob/master/lab2/README.md#the-go-language)
3. [Writing Go Code](https://github.com/uis-dat320-fall18/assignments/blob/master/lab2/README.md#writing-go-code)
4. [Installing and Running Go Code](https://github.com/uis-dat320-fall18/assignments/blob/master/lab2/README.md#installing-and-running-go-code)
5. [Go Exercises](https://github.com/uis-dat320-fall18/assignments/blob/master/lab2/README.md#go-exercises)
6. [Lab Approval](https://github.com/uis-dat320-fall18/assignments/blob/master/lab2/README.md#lab-approval)

## Introduction

In this lab you will be learning the Go programming language, which we will be using throughout
the course. Go is a very nice language in many ways. It provides built in primitives to design
concurrent programs using lightweight threads, called goroutines. In addition, these goroutines
can communicate using channels, instead of by means of shared memory which is the most
common approach in other languages such as C and Java. Go is also a very small language
in the sense that it has very few keywords, and as such it does not take too long to learn the
fundamentals of the language.

## The Go Language

There are many tutorials, books and other documentation on the Go language available on
the web. When searching the web, use the word golang instead of go. The main source of
information should preferably be [golang.org](http://golang.org/), which offers extensive documentation. A good
book for beginners is Miek Gieben’s [Learning Go](http://www.miek.nl/downloads/Go/Learning-Go-latest.pdf). Another great way to learn Go
is to complete the tour linked in the task below.

**Task: Start learning the basics of Go.** Complete the tour available on: [tour.golang.org](http://tour.golang.org/). You
should do at least the following exercises.

[Exercise: Loops and Functions](http://tour.golang.org/flowcontrol/8)

[Exercise: Slices](https://tour.golang.org/moretypes/14)

[Exercise: Maps](https://tour.golang.org/moretypes/19)

[Exercise: Errors](https://tour.golang.org/methods/19)

[Exercise: rot13Reader](https://tour.golang.org/methods/23)

Note that you can change the code inline in the browser and run the code to see the results.

Later, when you are more familiar with the language, the Go blog has a number of great
articles that describe important idioms in the language that may not be obvious; check it out
[here](http://blog.golang.org/index). And there is of course a number of youtube videos that may also be worthwhile to watch
to learn about specific topics, or the language proper.

**Go Language Questions**
```
1.) Which loop will repeat n times?
	a. for i := 1; i < n; i++ {}
	b. for i := 1; i <= n; i++ {}
	c. for i := 0; i <= n; i++ {}
	d. for i := 1; i == n; i++ {}

2.) What is the value of ok in the following example?

someMap := make(map[int]string)
someMap[0] = "String"
_, ok := someMap[0]

	a. "String"
	b. 'S'
	c. nil
	d. true

3.) Which function declaration is correct?

	a. func SampleFunction(x int, y string)
	b. func SampleFunction(int x, string y)
	c. method SampleFunction(x int, y string)
	d. int SampleFunction(int x, string y)

4.) Which function declaration will return an integer and an error?

	a. func SampleFunction(int x, error y)
	b. func SampleFunction(x int, y string) int, error
	c. int, error SampleFunction(x int, y string)
	d. func SampleFunction(x int, y string) (int, error)

5.) Which line passes a value x to a channel ch?
	a. ch.in(x)
	b. x <- ch
	c. ch <- x
	d. ch(x)
```

## Writing Go Code

There exist Go support for many editors, among others, vim and Eclipse. See [go-lang.cat-v.org/text-editors](http://go-lang.cat-v.org/text-editors/) for others.
Note that early versions of the Golang Eclipse plugin was pretty flaky, but this might have improved over the years, so feel free to try it out.
[Visual Studio Code](https://code.visualstudio.com/) is a great cross-platform editor with support for the Go language through [an extension](https://code.visualstudio.com/docs/languages/go). Another great editor with Go support is Sublime Text.

Here are some directions to get Sublime Text to work with Go. You will need to install three things.

`1.` [Sublime](https://www.sublimetext.com/3)

Download and extract the 64-bit Linux tarball.
```
	cd $HOME/Downloads
	tar -C <destination directory> -xvf <tarball file name>
```
Example:
```
	tar -C $HOME -xvf sublime_text_3_build_3103_x64.tar.bz2
```
Adding the location of the `sublime_text` executable to your path in the `.bashrc` file in your `$HOME` directory can make opening sublime easier.
```
	cd $HOME
	vi .bashrc
```
Add the following line:
```
	export PATH=$PATH:<directory of sublime_text>
```
Example:
```
	export PATH=$PATH:$HOME/sublime_text_3
```
Now you should be able to open Sublime Text just by typing
	`sublime_text`

`2.` [Package Control](https://packagecontrol.io/)

Just follow the instructions [here](https://packagecontrol.io/installation).

`3.` [GoSublime](https://github.com/DisposaBoy/GoSublime)

Just follow the instructions [here](https://github.com/DisposaBoy/GoSublime#installation).

Whichever editor you choose, it is highly recommended that you configure it to use the
[goimports](https://github.com/bradfitz/goimports) tool. This will reformat your code to follow the Go style, and make sure that all
the necessary import statements are inserted (so you don’t need to write the import statement
when you start using a new package.) The goimports tool is compatible with most editors, but
may require some configuration.

Note that editors may also be able to run your code within the editor itself, but it may
require some configuration. However, using the go tool from a terminal window (i.e. the
command line) is often times preferred.

## Installing and Running Go Code

The following step is optional, but included to show how to build a package (without a main
program). Why is that useful? Well, if you need to use an external package, i.e. a library, that
only provide API functions, and no main function, then this is the way to compile and install
it. So at this point, you can do the following to install the lab2 config package (again this step
is optional):

`go install github.com/uis-dat320-fall18/assignments/lab2/installation_task/config`

which will compile and install the files from the config directory as a library named
config.a in a subfolder of the pkg directory, which you can inspect as follows:

`ls -laR $GOPATH/pkg/linux_amd64/github.com/uis-dat320-fall18/assignments/lab2/installation_task`

Running go install compiles the package named on the command line (as shown above
with the config package.) However, the main.go file declares, package main in the first few
lines, and also contains a main() function. For this case we can do several things, all of which
produces runnable code, as outlined below:

`go install github.com/uis-dat320-fall18/assignments/lab2/installation_task`

which will compile and install a binary executable file called lab2 in $GOPATH/bin. Note
that lab2 is the name of the directory in which the main package is found, and is also given as
the last element of the path to the go install command. To run it simply do:

`$GOPATH/bin/installation_task`

Or you can also do:
```
cd $GOPATH/src/github.com/uis-dat320-fall18/assignments/lab2/installation_task
go build .
./installation_task
```
which will compile/build a binary executable file called lab2, save it in the current directory,
and finally run it. Or you can do:
```
cd $GOPATH/src/github.com/uis-dat320-fall18/assignments/lab2/installation_task
go run main.go
```
All these should result in the same output:
```
{1 hello}
{1 hello}
Text parsed correctly.
Gob parsed correctly.
{1 hello}
{1 hello}
```
These different approaches to compile/run a Go program may serve different uses depending
on what makes sense for different development approaches, and whether you need a binary
executable file and whether you need to install it in the default $GOPATH/bin.

Note that you may want to include $GOPATH/bin in your $PATH variable, which will allow
you to run the lab2 binary without prefixing it with a directory.

It is recommended that you use go install, since it will produce a binary that you can
distribute to other machines of the same architecture. But for quick testing go run is also very
convenient.

## Go Exercises

1. In the following, we will use **Task 1** as an example. Change directory to:
   `cd $GOPATH/src/github.com/uis-dat320-fall18/assignments/lab2` and confirm that the files
   for lab2 resides in that folder. They should, assuming that you ran the `go
   get` command in lab 1. The file `fib.go` contain the following skeleton code:

   ```go
   package lab1

   // Task 1: Fibonacci numbers
   //
   // fibonacci(n) returns nth Fibonacci number, and is defined by the
   // recurrence relation F_n = F_n-1 + F_n-2, with seed values F_0=0 and F_1=1.
   func fibonacci(n uint) uint {
	   return 0
   }
   ```

2. Implement the function body according to the specification so that all the
    tests in the corresponding `fib_test.go` file passes. The file looks like
    this:
    ```
    go
    package lab2

    import "testing"

    var fibonacciTests = []struct {
	    in, want uint
    }{
	    {0, 0},
	    {1, 1},
	    {2, 1},
	    {3, 2},
	    {4, 3},
	    {5, 5},
	    {6, 8},
	    {7, 13},
	    {8, 21},
	    {9, 34},
	    {10, 55},
	    {20, 6765},
    }

    func TestFibonacci(t *testing.T) {
	    for i, ft := range fibonacciTests {
		    out := fibonacci(ft.in)
		    if out != ft.want {
			    t.Errorf("fib test %d: got %d for input %d, want %d", i, out, ft.in, ft.want)
		    }
	    }
    }
    ```

3. If you run `go test` without any arguments, the tool will run all the tests
    found in files with name matching the pattern "*_test.go". You may only run
    a specific test by providing the `-run` flag to `go test`. For example, `go
    test -run TestFib` will only run the `TestFibonacci` test. Generally,
    running `go test -run regexp` will run only those tests matching the
    regular expression `regexp`.

4. You should ***not*** edit files or code that are marked with a `// DO NOT
    EDIT` comment. Please make separate `filename_test.go` files if you wish
    to write and run your own tests.

5. When you have completed a task and sufficiently many local tests pass, you
    may push your code to GitHub. This will trigger Autograder which will then
    run a separate test suite on your code.

6. Using the Fibonacci task (`fib.go`) as an example, use the following
    procedure to commit and push your changes to GitHub and Autograder:
    ```
    $ cd $GOPATH/src/github.com/uis-dat320-fall18/assignments/lab2
    $ git add fib.go
    $ git commit
    // This will open an editor for you to write a commit message
    // Use for example "Implemented Assignment 2"
    $ git push labs
    ```

7. Running the last command above will, due to an error on our part, result in
    Git printing an error message about a conflict between the `README.md` file
    in the `labs` repository and the `README.md` file in your `username-labs`
    repostitory. Here is how to fix it:

    ```
    $ git push labs
    ...
    ! [rejected]        master -> master (fetch first)
    error: failed to push some refs to 'git@github.com:uis-dat320-fall18/username-labs.git'
    ...
    $ git pull labs master
    ...
    Auto-merging README.md
    CONFLICT (add/add): Merge conflict in README.md
    Automatic merge failed; fix conflicts and then commit the result.
    ...
    $ cd $GOPATH/src/github.com/uis-dat320-fall18/assignments
    $ nano README.md
    // Remove everything in the file, then add for example "username-labs" to the file.
    // Save and exit.
    $ git add README.md
    $ git commit
    $ // Use the existing (merge) commit message. Save and exit.
    $ git push labs
    // Your push should now complete successfully.
    // You may check that your changes are reflected on GitHub through the GitHub web interface.
    ```

8. Autograder will now build and run a test suite on the code you submitted.
    You can check the output by going the [Autograder web
    interface](http://autograder.ux.uis.no/). The results (build log) should be
    available under "Individual - lab1". Note that the results shows output
    for all the tests in current lab assignment. You will want to focus on the
		output for the specific test results related to the task you're working on.

9. **Task 2:** Complete the task found in `stringer.go`. You may check your
    solution locally with the tests found in `stringer_test.go`.

10. **Task 3:** Complete the task found in `rot13.go`. You may check your
    solution locally with the tests found in `rot13_test.go`.

11. **Task 4:** Complete the task found in `scramble.go`. You may check your
    solution locally with the tests found in `scramble_test.go`

12. **Task 5:** Complete the task found in `errors.go`. You may check your
    solution locally with the tests found in `errors_test.go`.

13. **Task 6:** Complete the task found in `multiwriter.go`. You may check your
    solution locally with the tests found in `multiwriter_test.go`.

14. **Task 7:** Complete the task found in `password_manager.go`. The working
    program should be shown to the student assistants when approving the lab.

15. When you are finished with all the tasks for the current lab, and wish
    to submit, then first make sure you commit your changes and write only the
    following: `username labX submission` in the first line of the commit
    message, where you replace `username` with your GitHub username and `X`
    with the lab number. Your usage of slip days will be calculated based on
    when you pushed this commit to GitHub. If there are any issues you want us
    to pay attention to, please add those comments after an empty line in the
    commit message. If you later find a mistake and want to resubmit, please
    use `username labX resubmission` as the commit message.

16. Push your changes using `git push labs`. You should be able to view your
    results in the Autograder web interface as described earlier.

## Lab Approval

To have your lab assignment approved, you must come to the lab during lab hours
and present your solution. This lets you present the thought process behind your
solution, and gives us a more information for grading purposes. When you are
ready to show your solution, reach out to a member of the teaching staff.
It is expected that you can explain your code and show how it works.
You may show your solution on a lab workstation or your own
computer. The results from Autograder will also be taken into consideration
when approving a lab. At least 60% of the Autograder tests should pass for the
lab to be approved. A lab needs to be approved before Autograder will provide
feedback on the next lab assignment.

Also see the [Grading and Collaboration
Policy](https://github.com/uis-dat320-fall18/course-info/blob/master/policy.md)
document for additional information.

