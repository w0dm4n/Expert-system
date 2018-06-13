package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	EOF_TYPE   = "EOF"
	ERROR_TYPE = "error"
)

func errorChecker(err *error) {
	if *err != nil {
		panic(*err)
	}
}

func main() {
	defer func() {
		recover := recover()
		if recover != nil {
			err := recover.(error)
			if err != nil && typeOf(err) == ERROR_TYPE {
				fmt.Println("An error occured while trying to load datas:", err.Error())
			}
		}
	}()
	// log.SetFlags(0)
	// log.SetOutput(ioutil.Discard)
	var parser Parser
	parser.graph.Facts = make(map[string]*Fact)

	parser.graph.build()
	switch lenArgs := len(os.Args); lenArgs {
	case 1:
		fmt.Println("Please enter content and press CTRL-D")
		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString(0)
		if err != nil {
			if err.Error() == EOF_TYPE {
				parser.parseContent([]byte(data))
			} else {
				panic(err)
			}
		}
	case 2:
		data, err := ioutil.ReadFile(os.Args[1])
		errorChecker(&err)
		parser.parseContent(data)
	default:
		panic("Please give me a file or some entries in stdin :(")
	}
}
