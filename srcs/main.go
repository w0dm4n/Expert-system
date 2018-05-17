package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
			if err != nil {
				fmt.Println("An error occured while trying to load datas:", err.Error())
			}
		}
	}()

	switch lenArgs := len(os.Args); lenArgs {
	case 1:
		fmt.Println("Please enter content and press CTRL-D")
		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString(0)
		errorChecker(&err)
		parseContent(string(data))
	case 2:
		data, err := ioutil.ReadFile(os.Args[1])
		errorChecker(&err)
		parseContent(string(data))
	default:
		panic("Please give me a file or some entries in stdin :(")
	}
}
