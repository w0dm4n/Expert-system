package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
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
			fmt.Println(err)
			if err != nil && typeOf(err) == ERROR_TYPE {
				fmt.Println("An error occured while trying to load datas:", err.Error())
			}
		}
	}()
	var parser Parser
	parser.graph.Facts = make(map[string]*Fact)

	parser.graph.build()

	// check falgs
	var verbose bool
	var toDelete []int
	for i, arg := range os.Args {
		if arg == "-v" {
			verbose = true
			toDelete = append(toDelete, i)
		} else if arg == "-u" {
			parser.shouldRequestUndetermined = true
			toDelete = append(toDelete, i)
		}
	}

	// remove flags from args list
	count := 0
	for _, index := range toDelete {
		index = index - count
		os.Args = os.Args[:index+copy(os.Args[index:], os.Args[index+1:])]
		count++
	}
	if verbose == false {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	switch lenArgs := len(os.Args); lenArgs {
	case 1:
		fmt.Println("Please enter content and press CTRL-D")
		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString(0)
		if err != nil {
			if err.Error() == EOF_TYPE {
				fmt.Println()
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
