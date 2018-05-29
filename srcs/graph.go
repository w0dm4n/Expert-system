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
	setParentNode(Noder)
	apply() bool
}

type Rule struct {
	Type        string
	parentNodes []Noder
}

func (rule *Rule) getParentNodes() []Noder {
	return rule.parentNodes
}

func (rule *Rule) setParentNode(noder Noder) {
	rule.parentNodes = append(rule.parentNodes, noder)
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

func (fact *Fact) setParentNode(noder Noder) {
	fact.parentNodes = append(fact.parentNodes, noder)
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
	Facts     map[string]*Fact
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

func (graph *Graph) integrate(lhsNode *Node, op *BaseOperator, rhsNode *Node) {
	rootRule := &Rule{Type: op.Value}
	linked := graph.toNoder(lhsNode)
	invertLinked := graph.toNoder(rhsNode)
	println("linking")
	if rootRule.Type == "=>" {
		rootRule.parentNodes = append(rootRule.parentNodes, linked)
		graph.integrateNode(lhsNode, linked, true)
		graph.integrateNode(rhsNode, invertLinked, false)
	}
	if rootRule.Type == "<=>" {
		rootRule.parentNodes = append(rootRule.parentNodes, invertLinked)
		graph.integrateNode(rhsNode, invertLinked, true)
		graph.integrateNode(lhsNode, linked, false)
	}
}

// Integrate the current node and the left and right of the node
// noder is the equivalent of the node
func (graph *Graph) integrateNode(node *Node, noder Noder, isParent bool) {
	linkeds := make([]struct {
		*Node
		Noder
	}, 0)
	if node.Left != nil {
		linkeds = append(linkeds, struct {
			*Node
			Noder
		}{node.Left, graph.toNoder(node.Left)})
	}
	if node.Right != nil {
		linkeds = append(linkeds, struct {
			*Node
			Noder
		}{node.Right, graph.toNoder(node.Right)})
	}
	for _, linked := range linkeds {
		noder.setParentNode(linked.Noder)
		graph.integrateNode(linked.Node, linked.Noder, isParent)
	}
}

// 2 possible cases:
// we return an existing fact
// we return a new node which can be a new fact or just a rule
func (graph *Graph) toNoder(node *Node) (noder Noder) {
	if item, ok := graph.Facts[string(node.Value)]; ok {
		return item
	} else {
		if node.Value == '!' ||
			node.Value == ([]rune(SYMBOL_AND))[0] ||
			node.Value == ([]rune(SYMBOL_OR))[0] ||
			node.Value == ([]rune(SYMBOL_XOR))[0] {
			// got rule
			println("got rule")
			return &Rule{Type: string(node.Value)}
		} else {
			// got fact
			println("got fact")
			graph.Facts[string(node.Value)] = &Fact{Name: string(node.Value)}
			return graph.Facts[string(node.Value)]
		}
	}
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
