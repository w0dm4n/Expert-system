package main

import (
	"fmt"
	"log"
)

type Value int

const (
	False Value = iota
	True
	Undetermined
	DeadEnd
)

func (value Value) String() string {
	names := [...]string{
		"False",
		"True",
		"Undetermined",
		"DeadEnd"}

	if value < False || value > DeadEnd {
		return "Unknown"
	}
	return names[value]
}

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
	getChildNodes() []Noder
	setChildNode(Noder)
	apply(originsStack []*Fact, previous Noder, sameSide bool, visiteds map[Noder][]FactResult, requestingParents map[Noder][]FactRequest) Value
	name() string
}

type GraphNode struct {
	parentNodes []Noder
	childNodes  []Noder
}

func (node *GraphNode) getParentNodes() []Noder {
	return node.parentNodes
}

func (node *GraphNode) setParentNode(noder Noder) {
	node.parentNodes = append(node.parentNodes, noder)
}

func (node *GraphNode) getChildNodes() []Noder {
	return node.childNodes
}

func (node *GraphNode) setChildNode(noder Noder) {
	node.childNodes = append(node.childNodes, noder)
}

type Rule struct {
	Type string
	GraphNode
}

func (rule *Rule) name() string {
	return rule.Type
}

func (fact *Fact) name() string {
	return fact.Name
}

func opposite(value Value) Value {
	if value == True {
		return False
	} else if value == False {
		return True
	} else {
		return value
	}
}

func and(leftValue Value, rightValue Value) Value {
	if leftValue == False || rightValue == False {
		return False
	} else if leftValue == True && rightValue == True {
		return True
	} else if leftValue == Undetermined || rightValue == Undetermined {
		return Undetermined
	} else {
		return DeadEnd
	}
}

func or(leftValue Value, rightValue Value) Value {
	if leftValue == True || rightValue == True {
		return True
	} else if leftValue == False && rightValue == False {
		return False
	} else if leftValue == Undetermined || rightValue == Undetermined {
		return Undetermined
	} else {
		return DeadEnd
	}
}

// func or(leftValue Value, rightValue Value) Value {
// 	if leftValue == True || rightValue == True {
// 		return True
// 	} else {
// 		return Undefined
// 	}
// }

func xor(leftValue Value, rightValue Value) Value {
	if leftValue == True && rightValue == False || leftValue == False && rightValue == True {
		return True
	} else if leftValue != Undetermined && leftValue == rightValue {
		return False
	} else {
		return Undetermined
	}
}

// func xor(leftValue Value, rightValue Value) Value {
// 	if leftValue == True && rightValue == False || leftValue == False && rightValue == True {
// 		return True
// 	} else {
// 		return Undefined
// 	}
// }

func (rule *Rule) apply(originsStack []*Fact, previous Noder, sameSide bool, visiteds map[Noder][]FactResult, requestingParents map[Noder][]FactRequest) Value {

	var lastOrigin *Fact
	if len(originsStack) > 0 {
		lastOrigin = originsStack[len(originsStack)-1]
	}
	_ = lastOrigin

	if previous == rule.parentNodes[0] {
		log.Println("looking for child", rule.Type, "side", sameSide)
	} else {
		log.Println("looking for parent", rule.Type, "side", sameSide)
	}
	if rule.Type == "=>" {
		value := rule.parentNodes[0].apply(originsStack, rule, false, visiteds, requestingParents)
		if value == True {
			return True
		}
		return DeadEnd
	}

	if !sameSide {
		if rule.Type == "!" {
			return opposite(rule.parentNodes[0].apply(originsStack, rule, sameSide, visiteds, requestingParents))
		} else if rule.Type == "+" {
			return and(rule.parentNodes[0].apply(originsStack, rule, sameSide, visiteds, requestingParents), rule.parentNodes[1].apply(originsStack, rule, sameSide, visiteds, requestingParents))
		} else if rule.Type == "|" {
			return or(rule.parentNodes[0].apply(originsStack, rule, sameSide, visiteds, requestingParents), rule.parentNodes[1].apply(originsStack, rule, sameSide, visiteds, requestingParents))
		} else if rule.Type == "^" {
			return xor(rule.parentNodes[0].apply(originsStack, rule, sameSide, visiteds, requestingParents), rule.parentNodes[1].apply(originsStack, rule, sameSide, visiteds, requestingParents))
		} else {
			log.Println("woops unknown left side rule")
			return Undetermined
		}
	}

	// Need to explore children to infer a value..

	// We come from above
	// Need to know children combined value
	// children := rule.getChildNodes()
	// var leftValue Value
	// // if children[0] == lastOrigin {
	// // 	leftValue = DeadEnd
	// // } else {
	// leftValue = children[0].apply(originsStack, rule, sameSide, visiteds, requestingParents)
	// // }
	// var rightValue Value
	// if len(children) > 1 {
	// 	// if children[1] == lastOrigin {
	// 	// 	rightValue = DeadEnd
	// 	// } else if children[1] != nil {
	// 	rightValue = children[1].apply(originsStack, rule, sameSide, visiteds, requestingParents)
	// 	// }
	// }

	if previous == rule.parentNodes[0] {
		children := rule.getChildNodes()
		var leftValue Value
		// if children[0] == lastOrigin {
		// 	leftValue = DeadEnd
		// 	return DeadEnd
		// } else {
		leftValue = children[0].apply(originsStack, rule, sameSide, visiteds, requestingParents)
		// }
		var rightValue Value
		if len(children) > 1 {
			// if children[1] == lastOrigin {
			// 	rightValue = DeadEnd
			// 	return DeadEnd

			// } else if children[1] != nil {
			rightValue = children[1].apply(originsStack, rule, sameSide, visiteds, requestingParents)
			// }
		}
		if rule.Type == "!" {
			return opposite(leftValue)
		} else if rule.Type == "+" {
			// We come from above
			// Need to know children combined value
			return and(leftValue, rightValue)
		} else if rule.Type == "|" {
			return or(leftValue, rightValue)
		} else if rule.Type == "^" {
			return xor(leftValue, rightValue)
		}
	}

	// We come from below
	// need to know parent and other child value if parent is true or false
	var parentValue Value
	parentValue = rule.parentNodes[0].apply(originsStack, rule, sameSide, visiteds, requestingParents)

	log.Println("coming from below, checking other side")
	log.Println(rule.parentNodes[0].name(), "parent was", parentValue)
	if parentValue != True && parentValue != False {
		return parentValue
	}
	children := rule.getChildNodes()
	log.Println("applying", rule.Type)
	var other Noder
	for _, noder := range children {
		if noder != previous {
			other = noder
			break
		}
	}
	if rule.Type == "!" {
		return opposite(parentValue)
	} else if rule.Type == "+" {
		// check other child to detect inconsistency
		otherValue := other.apply(originsStack, rule, sameSide, visiteds, requestingParents)

		if parentValue == True {
			// just apply parent value directly
			return parentValue
		} else {
			if otherValue == True {
				return False
			} else {
				return Undetermined
			}
		}
	} else if rule.Type == "|" {
		// We come from below
		// need to know other child value
		otherValue := other.apply(originsStack, rule, sameSide, visiteds, requestingParents)
		log.Println("other value is", otherValue)
		if parentValue == True {
			if otherValue == False {
				return True
			} else {
				return Undetermined
			}
		} else {
			// if otherValue == False {
			return False
			// } else {
			// return Undefined
			// }
		}
	} else if rule.Type == "^" {
		// We come from below
		// need to know other child value
		otherValue := other.apply(originsStack, rule, sameSide, visiteds, requestingParents)
		if otherValue == False && parentValue == True {
			return True
		} else if otherValue == True && parentValue == True {
			return False
		} else if otherValue == True && parentValue == False {
			return True
		} else if otherValue == False && parentValue == False {
			return False
		} else {
			return Undetermined
		}
	} else {
		log.Println("woops unknown right side rule")
		return Undetermined
	}
}

type Fact struct {
	Name         string
	initialValue Value
	GraphNode
}

func bestValue(values []FactResult) Value {
	var gotTrue, gotFalse, gotUndefined, gotDeadEnd bool
	if len(values) == 0 {
		return Undetermined
	}
	for _, res := range values {
		if res.Value == True {
			gotTrue = true
		} else if res.Value == False {
			gotFalse = true
		} else if res.Value == Undetermined {
			gotUndefined = true
		} else if res.Value == DeadEnd {
			gotDeadEnd = true
		}
	}
	if gotTrue {
		return True
	} else if gotFalse {
		return False
	} else if gotUndefined {
		return Undetermined
	}
	_ = gotDeadEnd
	return DeadEnd
}

type FactResult struct {
	Value    Value
	Previous Noder
}

type FactRequest struct {
	origin   *Fact
	previous Noder
}

func stackPop(stack []*Fact) (res *Fact, resStack []*Fact) {
	if len(stack) > 0 {
		res, resStack = stack[len(stack)-1], stack[:len(stack)-1]
	} else {
		res, resStack = nil, stack
	}
	return
}

func stackContains(stack []*Fact, fact *Fact) bool {
	for _, item := range stack {
		if item == fact {
			return true
		}
	}
	return false
}

func (fact *Fact) apply(originsStack []*Fact, previous Noder, sameSide bool, visiteds map[Noder][]FactResult, requestingParents map[Noder][]FactRequest) Value {

	entryVisiteds := make(map[Noder][]FactResult)
	entryRequestingParents := make(map[Noder][]FactRequest)
	_, _ = entryVisiteds, entryRequestingParents

	var lastOrigin *Fact
	if len(originsStack) > 0 {
		lastOrigin = originsStack[len(originsStack)-1]
	}

	// if visited already, return its value
	if _, ok := visiteds[fact]; ok && !sameSide {
		log.Println("looking for existing", fact.Name)
		log.Println("------", fact.Name, "is", visiteds[fact])
		value := bestValue(visiteds[fact])
		if value == True || value == False {
			return value
		}
		// return bestValue(visiteds[fact])
	}
	log.Println("looking for", fact.Name)

	// var lastOrigin *Fact
	// if len(originsStack) > 0 {
	// 	lastOrigin = originsStack[len(originsStack)-1]
	// }
	// if lastOrigin == fact && sameSide {
	// 	log.Println("last origin was the same on the same side")
	// 	// res := FactResult{Value: DeadEnd, Previous: previous}
	// 	// log.Println(fact.Name, "was", visiteds[fact], "and got", res)
	// 	// visiteds[fact] = append(visiteds[fact], res)
	// 	return bestValue(visiteds[fact])
	// }

	// log.Println("make sure we don't already come from that parent")
	if parents, ok := requestingParents[fact]; ok {
		for _, request := range parents {
			if request.previous == previous && request.origin == lastOrigin {
				var resValue Value
				// if fact.initialValue == True {
				// 	resValue = True
				// } else {
				// 	resValue = origin.initialValue
				// }

				resValue = DeadEnd
				// resValue = bestValue(visiteds[fact])

				log.Println("parent visited already,", resValue)
				var res FactResult
				// res = FactResult{Value: origin.initialValue, Previous: previous}
				// log.Println("[origin]", origin.Name, "was", visiteds[origin], "and got", res)
				// visiteds[origin] = append(visiteds[origin], res)

				res = FactResult{Value: resValue, Previous: previous}
				log.Println(fact.Name, "was", visiteds[fact], "and got", res)
				visiteds[fact] = append(visiteds[fact], res)
				return resValue
				// return False
				// return DeadEnd
			}
		}
	}
	req := FactRequest{origin: lastOrigin, previous: previous}
	if req.origin != nil {
		log.Println("requesting parent is", req.origin.Name)
		// if req.origin == fact {
		// 	log.Println("same as current, returning", bestValue(visiteds[fact]))
		// 	return bestValue(visiteds[fact])
		// }
	}
	requestingParents[fact] = append(requestingParents[fact], req)

	if len(fact.parentNodes) == 0 /* || (fact.initialValue == True && !sameSide)*/ {
		log.Println("applying default value")
		res := FactResult{Value: fact.initialValue, Previous: previous}
		log.Println(fact.Name, "was", visiteds[fact], "and got", res)
		visiteds[fact] = append(visiteds[fact], res)
		// log.Println(fact.Name, "is", visiteds[fact])
		return bestValue(visiteds[fact])
	}

	// if lastOrigin == fact && sameSide {
	// 	log.Println("last origin was the same on the same side")
	// 	// res := FactResult{Value: DeadEnd, Previous: previous}
	// 	// log.Println(fact.Name, "was", visiteds[fact], "and got", res)
	// 	// visiteds[fact] = append(visiteds[fact], res)
	// 	return DeadEnd
	// }

	// potentialsValues := make([]Value, len(fact.parentNodes))
	for _, v := range fact.parentNodes {
		// if previous == v && sameSide {
		// 	// parent is asking for the value, so we can't get one this way
		// 	// potentialsValues[i] = Undefined
		// 	log.Println("parent is asking for the value")
		// 	res := FactResult{Value: DeadEnd, Previous: previous}
		// 	log.Println(fact.Name, "was", visiteds[fact], "and got", res)
		// 	visiteds[fact] = append(visiteds[fact], res)
		// } else {
		// potentialsValues[i] = v.apply(fact, true, visiteds)
		// visiteds[fact] = append(visiteds[fact], potentialsValues[i])
		originsStack = append(originsStack, fact)
		res := FactResult{Value: v.apply(originsStack, fact, true, visiteds, requestingParents), Previous: previous}
		_, originsStack = stackPop(originsStack)
		log.Println(fact.Name, "was", visiteds[fact], "and got", res)
		visiteds[fact] = append(visiteds[fact], res)
		// }
	}

	// make sure we have a definitive true or false
	var gotTrue, gotFalse, gotUndefined, gotDeadEnd bool
	for _, res := range visiteds[fact] {
		if res.Value == True {
			gotTrue = true
		} else if res.Value == False {
			gotFalse = true
		} else if res.Value == Undetermined {
			gotUndefined = true
		} else if res.Value == DeadEnd {
			gotDeadEnd = true
		}
	}
	if gotTrue && gotFalse {
		panic(fmt.Sprint("opposite conditions on ", fact.Name))
	}
	//  else if gotTrue {
	// 	visiteds[fact] = append(visiteds[fact], True)
	// } else if gotFalse {
	// 	visiteds[fact] = append(visiteds[fact], False)
	// } else if gotUndefined || gotDeadEnd {
	// 	visiteds[fact] = append(visiteds[fact], Undefined)
	// } else {
	// 	visiteds[fact] = append(visiteds[fact], False)
	// }
	_ = gotUndefined
	_ = gotDeadEnd

	log.Println("======", fact.Name, "is", visiteds[fact])
	value := bestValue(visiteds[fact])
	// if lastOrigin == nil && value == DeadEnd {
	// 	// if !stackContains(originsStack, fact) && value == DeadEnd {
	// 	log.Println(fact.Name, "is dead end and final, need to add its initial value and retry")
	// 	res := FactResult{Value: fact.initialValue, Previous: previous}
	// 	entryVisiteds[fact] = append(entryVisiteds[fact], res)
	// 	log.Println(fact.Name, "was", visiteds[fact], "and got", res)
	// 	originsStack = append(originsStack, fact)
	// 	value = fact.apply(originsStack, previous, true, entryVisiteds, entryRequestingParents)
	// 	_, originsStack = stackPop(originsStack)
	// }
	// if stackContains(originsStack, fact) && value == DeadEnd {
	// 	log.Println(fact.Name, "is DeadEnd and final, need to add its initial value and retry")
	// 	res := FactResult{Value: fact.initialValue, Previous: previous}
	// 	entryVisiteds[fact] = append(entryVisiteds[fact], res)
	// 	log.Println(fact.Name, "was", visiteds[fact], "and got", res)
	// 	originsStack = append(originsStack, fact)
	// 	value = fact.apply(originsStack, previous, true, entryVisiteds, entryRequestingParents)
	// 	_, originsStack = stackPop(originsStack)
	// }
	// if value == DeadEnd {
	// 	log.Println("fact is dead end, need to add its initial value")
	// 	res := FactResult{Value: fact.initialValue, Previous: previous}
	// 	log.Println(fact.Name, "was", visiteds[fact], "and got", res)
	// 	visiteds[fact] = append(visiteds[fact], res)
	// 	value = res.Value
	// }
	if value == DeadEnd /*&& len(originsStack) == 0*/ {
		value = False
	}
	return value
}

func printRulesUntilFact(noder Noder) {
	for _, parent := range noder.getParentNodes() {
		if fact, ok := parent.(*Fact); ok {
			// log.Println("converted into fact", fact)
			// printRules(fact)
			log.Print(fact.Name)
			_ = fact
		} else if rule, ok := parent.(*Rule); ok {
			// log.Println("converted into rule", rule)
			log.Print(rule.Type)
			printRulesUntilFact(rule)
			if rule.Type == "=>" {
				log.Println()
			}
		}
	}
}

func (fact *Fact) printRulesUntilFact() {
	log.Print(fact.Name)
	log.Println(" has", len(fact.getParentNodes()), "parent nodes")
	printRulesUntilFact(fact)
	log.Println()
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
	// log.Println("linking")
	if rootRule.Type == "=>" {
		rootRule.setParentNode(linked)
		linked.setChildNode(rootRule)
		invertLinked.setParentNode(rootRule)
		rootRule.setChildNode(invertLinked)
		graph.integrateNode(lhsNode, linked, true)
		graph.integrateNode(rhsNode, invertLinked, false)
	}
	if rootRule.Type == "<=>" {
		rootRule.setParentNode(invertLinked)
		invertLinked.setChildNode(rootRule)
		linked.setParentNode(rootRule)
		rootRule.setChildNode(linked)
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
		if isParent {
			noder.setParentNode(linked.Noder)
			linked.Noder.setChildNode(noder)
		} else {
			linked.Noder.setParentNode(noder)
			noder.setChildNode(linked.Noder)
		}
		graph.integrateNode(linked.Node, linked.Noder, isParent)
	}
}

// 2 possible cases:
// we return an existing fact
// we return a new node which can be a new fact or just a rule
func (graph *Graph) toNoder(node *Node) (noder Noder) {
	if item, ok := graph.Facts[string(node.Value)]; ok {
		// log.Println("got existing fact")
		return item
	} else {
		if node.Value == '!' ||
			node.Value == ([]rune(SYMBOL_AND))[0] ||
			node.Value == ([]rune(SYMBOL_OR))[0] ||
			node.Value == ([]rune(SYMBOL_XOR))[0] {
			// got rule
			// log.Println("got rule")
			return &Rule{Type: string(node.Value)}
		} else {
			// got fact
			// log.Println("got new fact", string(node.Value))
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
	if fact, ok := graph.Facts[string(operand)]; ok {
		fact.initialValue = True
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
