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
		content = content[(index + 1):len(content)]
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
}

func (parser *Parser) getQueryResult(content string, l int) {
	operands := []byte(strings.Trim(content, " "))
	for _, elem := range operands {
		operand := parser.graph.getOperand(rune(elem))
		if operand != nil {
			fmt.Printf("%s is %t\n", string(operand.Value), operand.Active)
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
			} else if strings.Index(elem, INITIAL_QUERIES) != -1 && strings.Index(elem, INITIAL_QUERIES) == 0 {

				// execute operations here
				parser.getQueryResult(elem[1:len(elem)], l)
			} else {
				panic(fmt.Sprintf("%s %d: %s", "Bad syntax on line", l, "No operator found"))
			}
		}
		l++
	}
}
