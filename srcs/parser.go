package main

import (
	"fmt"
	"strings"
)

const (
	INITIAL_FACTS     = "="
	INITIAL_QUERIES   = "?"
	NEGATIVE_OPERATOR = "!"
	PARENTHESIS_START = "("
	PARENTHESIS_END   = ")"
)

type Parser struct {
	graph Graph
}

type Node struct {
	Value  rune
	Parent *Node
	Left   *Node
	Right  *Node
}

func (parser *Parser) removeComment(content string) string {
	comment := strings.Index(content, "#")
	if comment != -1 {
		content = content[0:comment]
	}
	return content
}

func (parser *Parser) trimOperand(content string) string {
	var index = strings.Index(content, NEGATIVE_OPERATOR)
	if index != -1 {
		content = parser.trimOperand(content[(index + 1):len(content)])
	}
	index = strings.Index(content, PARENTHESIS_START)
	if index != -1 {
		content = content[(index + 1):len(content)]
	}
	index = strings.Index(content, PARENTHESIS_END)
	if index != -1 {
		content = content[0:index]
	}
	return content
}

func (parser *Parser) parseOperands(line []string) {
	for _, content := range line {
		content = strings.TrimSpace(content)
		content = parser.trimOperand(strings.ToUpper(content))

		if len(content) == 1 {
			operand := rune(content[0])

			if operand >= 'A' && operand <= 'Z' {
				if !parser.graph.operandExist(operand) {
					parser.graph.addOperand(operand)
				}
			}
		}
	}
}

func (parser *Parser) getOperator(content string) *BaseOperator {
	for _, elem := range parser.graph.Operators {
		index := strings.Index(content, elem.Value)
		if index != -1 {
			return &elem
		}
	}
	return nil
}

func (parser *Parser) activeOperands(content string, l int) {
	operands := []rune(strings.Trim(content, " "))
	for _, operand := range operands {
		// fmt.Println(operand)
		if operand >= 'A' && operand <= 'Z' {
			if parser.graph.operandExist(operand) {
				parser.graph.activeOperand(operand)
			} else {
				panic(fmt.Sprintf("%s %d: %s (%s)", "Bad syntax on line", l, "Operand do not exist or not used", string(operand)))
			}
		} else {
			panic(fmt.Sprintf("%s %d: %s", "Bad syntax on line", l, "Invalid operand syntax on initial fact"))
		}
	}
}

func (parser *Parser) getSymbol(content string) *BaseSymbol {
	for _, elem := range parser.graph.Symbols {
		if elem.Value == content {
			return &elem
		}
	}
	return nil
}

func (parser *Parser) getOperand(content string) *Operand {
	operandValue := []rune(parser.trimOperand(content))[0]
	if operandValue >= 'A' && operandValue <= 'Z' {
		operand := Operand{operandValue, true}
		if strings.Contains(content, NEGATIVE_OPERATOR) {
			operand.Active = false
		}
		return &operand
	}
	return nil
}

func (parser *Parser) getOperandConcerned(operands []Operand, content []string, i int) []Operand {
	var concerned []Operand
	concerned = append(concerned, operands[i-1])
	concerned = append(concerned, operands[i])
	return concerned
}

func (parser *Parser) newOperation(conditional, affected string, operator *BaseOperator) {

	conditionalContent := strings.Split(conditional, " ")
	var operands []Operand
	var symbols []Symbol
	var inParenthesis = false

	_ = inParenthesis
	for _, elem := range conditionalContent {

		operand := parser.getOperand(elem)
		if operand != nil {
			operands = append(operands, *operand)
		}
	}

	symbolCount := 0
	for _, elem := range conditionalContent {
		if strings.Contains(elem, PARENTHESIS_START) {
			inParenthesis = true
		}
		symbolBase := parser.getSymbol(elem)
		if symbolBase != nil {
			symbolCount++ // until next symbol
			symbol := Symbol{symbolBase.Value, parser.getOperandConcerned(operands, conditionalContent, symbolCount), inParenthesis}
			symbols = append(symbols, symbol)
		}
		if strings.Contains(elem, PARENTHESIS_END) {
			inParenthesis = false
		}
	}

	// C | !X + (B + X | (F | !X))
	for _, elem := range symbols {
		fmt.Printf("%s (%t) %s %s (%t) (Parenthese %t)\n", string(elem.OperandsAffected[0].Value), elem.OperandsAffected[0].Active, elem.Value, string(elem.OperandsAffected[1].Value), elem.OperandsAffected[1].Active, elem.inParenthesis)
	}

	// fmt.Println(symbols, "\n", operands)
	// for _, elem := range operands {
	// 	fmt.Printf("%s %t\n", string(elem.Value), elem.Active)
	// }

	// parenthesis := strings.Split(conditional, PARENTHESIS_START)
	// fmt.Println(parenthesis)
	fmt.Printf("Conditional: %s, Operator: %s, Affected: %s\n", conditional, operator.Value, affected)

	// rule creation into graph

	// transform conditional into a tree
	// starting node is always the operator
	Rule := &Rule{Type: operator.Value}

	_ = Rule

	lhsRawNodes, _ := arrangeOperations(conditional)
	rhsRawNodes, _ := arrangeOperations(affected)

	fmt.Println("actual tree")
	rhsRawNodes.print(1)
	fmt.Println(operator.Value)
	lhsRawNodes.print(1)

	lhsRawNodes = optimiseTree(lhsRawNodes)
	lhsRawNodes = optimiseTree(lhsRawNodes)
	rhsRawNodes = optimiseTree(rhsRawNodes)
	rhsRawNodes = optimiseTree(rhsRawNodes)

	fmt.Println("final optimized")
	rhsRawNodes.print(1)
	fmt.Println(operator.Value)
	lhsRawNodes.print(1)

	// conversion of binary tree nodes into graph nodes
	// the graph has to know on which side it is from the operator
	parser.graph.integrate(lhsRawNodes, operator, rhsRawNodes)
}

func optimiseTree(node *Node) (root *Node) {
	root = node

	// !! => nothing
	if node.Value == '!' && node.Left != nil && node.Left.Value == '!' {
		if root.Parent != nil {
			if root.Parent.Left == node {
				root.Parent.Left = node.Left.Left
				root.Parent.Left.Parent = root.Parent.Left
			} else {
				root.Parent.Right = node.Left.Left
				root.Parent.Right.Parent = root.Parent.Right
			}
		}
		root = node.Left.Left
		node = root
	}

	// !(A + B) => !A | !B
	if node.Value == '!' && node.Left != nil && node.Left.Value == '+' {
		_, child := remove(node)
		root = child
		root.Value = '|'
		insertBetween(root, root.Left, '!')
		insertBetween(root, root.Right, '!')
		node = root
	}

	// !(A | B) => !A + !B
	if node.Value == '!' && node.Left != nil && node.Left.Value == '|' {
		_, child := remove(node)
		root = child
		root.Value = '+'
		insertBetween(root, root.Left, '!')
		insertBetween(root, root.Right, '!')
		node = root
	}

	// !(A + !B | !A + B) => !(A + !B) | !(!A + B) => (!A | !!B) | (!!A | B) => (!A | B) | (A | !B)

	// A ^ B => A + !B | !A + B
	if node.Value == '^' {
		root = &Node{Value: '|', Left: node, Parent: node.Parent}
		if root.Parent != nil {
			if root.Parent.Left == node {
				root.Parent.Left = root
			} else {
				root.Parent.Right = root
			}
		}

		root.Left.Parent = root
		root.Left.Value = '+'

		root.Right = copyTree(root.Left)
		insertBetween(root.Left, root.Left.Left, '!')
		insertBetween(root.Right, root.Right.Right, '!')

		node = root
	}

	if root.Left != nil {
		root.Left = optimiseTree(root.Left)
	}
	if root.Right != nil {
		root.Right = optimiseTree(root.Right)
	}
	return root
}

func copyTree(node *Node) (copy *Node) {
	if node == nil {
		copy = node
		return
	}

	//create new node and make it a copy of node pointed by root
	copy = &Node{Value: node.Value, Left: copyTree(node.Left), Right: copyTree(node.Right), Parent: node.Parent}
	return
}

func remove(node *Node) (root *Node, child *Node) {
	child = node.Left
	if node.Parent != nil {
		if node.Parent.Left == node {
			node.Left.Parent = node.Parent
			node.Parent.Left = node.Left
		} else {
			node.Left.Parent = node.Parent
			node.Parent.Right = node.Left
		}
		root = node.Parent
	} else {
		root = node.Left
	}
	return
}

func insertBetween(parentNode *Node, childNode *Node, value rune) (root *Node) {
	root = parentNode
	newNode := &Node{Value: value, Left: childNode, Parent: parentNode}
	if parentNode == nil {
		root = newNode
	} else if parentNode.Left == childNode {
		parentNode.Left = newNode
	} else if parentNode.Right == childNode {
		parentNode.Right = newNode
	}
	childNode.Parent = newNode
	return
}

func (node *Node) print(level int) {
	if node.Right != nil {
		// fmt.Print("r ")
		node.Right.print(level + 1)
	}
	fmt.Print(level)
	for i := 0; i < level; i++ {
		fmt.Printf("  ")
	}
	fmt.Println(string(node.Value))
	if node.Left != nil {
		// fmt.Print("l ")
		node.Left.print(level + 1)
	}
}

var prios = map[rune]int{
	'(': 1,
	')': 2,
	'!': 3,
	([]rune(SYMBOL_AND))[0]: 4,
	([]rune(SYMBOL_OR))[0]:  5,
	([]rune(SYMBOL_XOR))[0]: 6,
}

func (node *Node) insert(currentNode *Node, value rune) (root *Node, inserted *Node) {
	// in case we come from root brackets (i.e. prev was set to null to force absolute priority)
	if currentNode == nil {
		currentNode = &Node{Value: value, Left: node}
		node.Parent = currentNode
		inserted = currentNode
		root = inserted
		return
	}

	root = node
	if prios[value] < prios[currentNode.Value] {
		if currentNode.Left == nil {
			currentNode.Left = &Node{Value: value, Parent: currentNode}
			inserted = currentNode.Left
		} else if currentNode.Right == nil {
			currentNode.Right = &Node{Value: value, Parent: currentNode}
			inserted = currentNode.Right
		} else {
			intermediateNode := &Node{Value: value, Parent: currentNode, Left: currentNode.Right}
			currentNode.Right.Parent = intermediateNode
			currentNode.Right = intermediateNode
			inserted = intermediateNode
		}
	} else {
		if currentNode.Parent != nil {
			root, inserted = node.insert(currentNode.Parent, value)
		} else {
			currentNode.Parent = &Node{Value: value, Left: currentNode}
			inserted = currentNode.Parent
			root = inserted
		}
	}

	return
}

func (node *Node) insertNode(currentNode *Node, incomingNode *Node) (root *Node, inserted *Node) {
	root = node
	if currentNode.Left == nil {
		currentNode.Left = incomingNode
		incomingNode.Parent = currentNode
		inserted = currentNode.Left
	} else if currentNode.Right == nil {
		currentNode.Right = incomingNode
		incomingNode.Parent = currentNode
		inserted = currentNode.Right
	}

	return
}

func arrangeOperations(operations string) (res *Node, length int) {

	var root *Node
	prev := root
	skip := 0

	// fmt.Println("arranging", operations)
	for pos, char := range []rune(operations) {
		// prev := root
		if skip > 0 {
			skip--
			continue
		}
		if char == ' ' || char == '\r' {
			continue
		}
		// fmt.Printf("character %c starts at byte position %d\n", char, pos)

		switch char {
		case '(':
			// fmt.Println("opening bracket")
			var innerOps *Node
			innerOps, length = arrangeOperations(operations[pos+1:])
			skip = length

			// fmt.Println("got back from resursive with")
			// innerOps.print(0)
			// fmt.Println("length was", skip)

			if root == nil {
				root = innerOps
				prev = root
			} else {
				// fmt.Println("[bracket] inserting", string(innerOps.Value), "on", string(prev.Value))
				root, prev = root.insertNode(prev, innerOps)
			}
			prev = prev.Parent
		case ')':
			// fmt.Println("closing bracket at", pos)
			return root, pos + 1
		default:
			if root == nil {
				root = &Node{Value: char}
				prev = root
			} else {
				// if prev != nil {
				// 	fmt.Println("inserting", string(char), "on", string(prev.Value))
				// } else {
				// 	fmt.Println("inserting", string(char), "on", nil)
				// }
				root, prev = root.insert(prev, char)
				for prev.Value == '!' && prev.Left != nil && prev.Left.Value == '!' {
					prev = prev.Left
				}
			}
		}
		// fmt.Println("current tree")
		// root.print(0)
		length++
	}

	res = root
	return
}

func (parser *Parser) getQueryResult(content string, l int) {
	operands := []byte(strings.Trim(content, " "))
	for _, elem := range operands {
		operand := parser.graph.getOperand(rune(elem))
		if operand != nil {
			fmt.Printf("%s is %t\n", string(operand.Value), operand.Active)
			if fact, ok := parser.graph.Facts[string(elem)]; ok {
				// fmt.Println("inferring value of", fact)
				fact.printRulesUntilFact()
				fmt.Println("value of", fact.Name, "was", fact.initialValue)
				visiteds := make(map[Noder][]FactResult)
				for _, fact := range parser.graph.Facts {
					if fact.initialValue == True {
						res := FactResult{Value: True, Previous: nil}
						visiteds[fact] = append(visiteds[fact], res)
					}
				}
				requestingParents := make(map[Noder][]FactRequest)
				originsStack := []*Fact{}
				value := fact.apply(originsStack, nil, true, visiteds, requestingParents)
				if value == DeadEnd {
					value = fact.initialValue
				}
				fmt.Println("value of", fact.Name, "is now", value)
			} else {
				fmt.Println("no fact registered for", string(elem))
			}
		} else {
			panic(fmt.Sprintf("%s %d: %s (%s)", "Bad syntax on line", l, "Invalid operand on query (do not exist or not used)", string(operand.Value)))
		}
	}
}

func (parser *Parser) parseContent(bytes []byte) {

	defer func() {
		recover := recover()
		if recover != nil {
			err := recover.(string)
			fmt.Println(err)
		}
	}()

	lines := strings.Split(string(bytes), "\n")
	l := 1
	for _, elem := range lines {
		elem = parser.removeComment(elem)
		elem = strings.TrimSpace(elem)
		if len(elem) > 0 {
			parser.parseOperands(strings.Split(elem, " "))
			operator := parser.getOperator(elem)
			if operator != nil {

				indexOperator := strings.Index(elem, operator.Value)
				operandsConditional := strings.Trim(elem[0:indexOperator], " ")
				operandsAffected := strings.Trim(elem[(indexOperator+len(operator.Value)):len(elem)], " ")

				parser.newOperation(operandsConditional, operandsAffected, operator)
			} else if strings.Index(elem, INITIAL_FACTS) != -1 && strings.Index(elem, INITIAL_FACTS) == 0 {
				parser.activeOperands(elem[1:len(elem)], l)
			} else if strings.Index(elem, INITIAL_QUERIES) == 0 {
				// execute operations here
				parser.getQueryResult(elem[1:len(elem)], l)
			} else {
				panic(fmt.Sprintf("%s %d: %s", "Bad syntax on line", l, "No operator found"))
			}
		}
		l++
	}

	// parser.graph.printConnections()
}
