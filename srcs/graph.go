package main

type Operand struct {
	Value  rune
	Active bool
}

type Symbol struct {
	Value            string
	OperandsAffected []Operand
}

type Graph struct {
	Operands  []Operand
	Operators [2]BaseOperator
	Symbols   [3]BaseSymbol
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
