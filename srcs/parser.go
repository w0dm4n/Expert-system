package main

import "strings"

type LinkType struct {
	Type     byte
	Constant []Constant
}

type Constant struct {
	ConstantName byte
	Linked       []LinkType
}

type Graph struct {
	Constants []Constant
}

func parseContent(content string) {
	lines := strings.Split(content, "\n")
	_ = lines
}
