package main

import "fmt"

type Operand struct {
	Value  rune
	Active bool
}

type Symbol struct {
	Value            string
	OperandsAffected []Operand
	inParenthesis    bool
}

type Noder interface {
	getParentNodes() []Noder
	apply() bool
}

type Rule struct {
	Type        string
	parentNodes []Noder
}

func (rule *Rule) getParentNodes() []Noder {
	return rule.parentNodes
}

func (rule *Rule) apply() bool {
	potentialsValues := make([]bool, len(rule.getParentNodes()))
	for i, v := range rule.parentNodes {
		potentialsValues[i] = v.apply()
	}
	// need to return a definitive value or undetermined here
	return potentialsValues[0]
}

type Fact struct {
	Name         string
	initialValue bool
	parentNodes  []Noder
}

func (fact *Fact) getParentNodes() []Noder {
	return fact.parentNodes
}

func (fact *Fact) apply() bool {
	if len(fact.parentNodes) == 0 {
		return fact.initialValue
	} else {
		potentialsValues := make([]bool, len(fact.parentNodes))
		for i, v := range fact.parentNodes {
			potentialsValues[i] = v.apply()
		}
		// need to return a definitive value or undetermined here
		return potentialsValues[0]
	}
}

type Graph struct {
	Facts     map[string]Fact
	Operands  []Operand
	Operators [2]BaseOperator
	Symbols   [3]BaseSymbol
}

func (graph *Graph) printConnections() {
	fmt.Printf("%+v", graph.Facts)
}

func (graph *Graph) print() {
	fmt.Printf("%+v", graph)
}

func (graph *Graph) addOperand(operand rune) {
	graph.Operands = append(graph.Operands, Operand{operand, false})
}

func (graph *Graph) operandExist(operand rune) bool {
	for _, elem := range graph.Operands {
		if elem.Value == operand {
			return true
		}
	}
	return false
}

func (graph *Graph) build() {
	graph.Operators = [...]BaseOperator{BaseOperator{OPERATOR_IF_ONLY}, BaseOperator{OPERATOR_IMPLIES}}
	graph.Symbols = [...]BaseSymbol{BaseSymbol{SYMBOL_AND}, BaseSymbol{SYMBOL_OR}, BaseSymbol{SYMBOL_XOR}}
}

func (graph *Graph) activeOperand(operand rune) {
	for _, elem := range graph.Operands {
		if elem.Value == operand {
			elem.Active = true
			break
		}
	}
}

func (graph *Graph) getOperand(operand rune) *Operand {
	for _, elem := range graph.Operands {
		if elem.Value == operand {
			return &elem
		}
	}
	return nil
}
