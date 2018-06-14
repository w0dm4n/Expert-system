package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestErrorFalseChain(t *testing.T) {
// 	var parser Parser
// 	parser.graph.Facts = make(map[string]*Fact)

// 	parser.graph.build()
// 	data, err := ioutil.ReadFile("../tests/testErrorFalseChain")
// 	errorChecker(&err)
// 	t.Log(string(data))

// 	assert.PanicsWithValue(t, "opposite conditions on K", func() { parser.parseContent(data) }, "failed to panic correctly")
// }

func TestMain(m *testing.M) {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)
	code := m.Run()
	os.Exit(code)
}

func testFileError(t *testing.T, file string, expectedError string) {
	var parser Parser
	assert.PanicsWithValue(t, expectedError, func() { parse(t, &parser, file) }, "failed to panic correctly "+expectedError)
}

// func TestErrorFalse(t *testing.T) {
// 	var parser Parser
// 	parser.graph.Facts = make(map[string]*Fact)

// 	parser.graph.build()
// 	data, err := ioutil.ReadFile("../tests/testErrorFalse")
// 	errorChecker(&err)
// 	t.Log(string(data))

// 	assert.PanicsWithValue(t, "opposite conditions on K", func() { parser.parseContent(data) }, "failed to panic correctly")
// }

type StdRedirector struct {
	rescueStdout *os.File
	r            *os.File
	w            *os.File
}

func (stdRedirector *StdRedirector) startRedirect() {
	stdRedirector.rescueStdout = os.Stdout
	stdRedirector.r, stdRedirector.w, _ = os.Pipe()
	os.Stdout = stdRedirector.w
}

func (stdRedirector *StdRedirector) endRedirect() string {
	stdRedirector.w.Close()
	out, _ := ioutil.ReadAll(stdRedirector.r)
	outStr := string(out)
	os.Stdout = stdRedirector.rescueStdout
	return outStr
}

var redirector StdRedirector

func parse(t *testing.T, parser *Parser, file string) {
	parser.graph.Facts = make(map[string]*Fact)

	parser.graph.build()
	data, err := ioutil.ReadFile(file)
	errorChecker(&err)
	t.Log(file)
	t.Log(string(data))
	parser.parseContent(data)
}

// func TestFalseToFalse(t *testing.T) {
// 	redirector.startRedirect()

// 	var parser Parser
// 	parse(t, &parser, "../tests/testFalseToFalse")

// 	outStr := redirector.endRedirect()
// 	t.Log(outStr)
// 	if !strings.Contains(outStr, "value of K is now False") {
// 		t.Error("K not false")
// 	}
// }

// func TestNotTrueToNotTrue(t *testing.T) {
// 	redirector.startRedirect()

// 	var parser Parser
// 	parse(t, &parser, "../tests/testNotTrueToNotTrue")

// 	outStr := redirector.endRedirect()
// 	t.Log(outStr)
// 	if !strings.Contains(outStr, "value of K is now False") {
// 		t.Error("K not false")
// 	}
// }

// func TestTrueToTrue(t *testing.T) {
// 	redirector.startRedirect()

// 	var parser Parser
// 	parse(t, &parser, "../tests/testTrueToTrue")

// 	outStr := redirector.endRedirect()
// 	t.Log(outStr)
// 	if !strings.Contains(outStr, "value of K is now True") {
// 		t.Error("K not True")
// 	}
// }

func testFile(t *testing.T, file string, expectedResult string) {
	redirector.startRedirect()

	var parser Parser
	parse(t, &parser, file)

	outStr := redirector.endRedirect()
	t.Log(outStr)
	if !strings.Contains(outStr, expectedResult) {
		t.Error("NOT", expectedResult)
	}
}

func testFileResults(t *testing.T, file string, expectedResults []string) {
	redirector.startRedirect()

	var parser Parser
	parse(t, &parser, file)

	outStr := redirector.endRedirect()
	t.Log(outStr)
	for _, expRes := range expectedResults {
		if !strings.Contains(outStr, expRes) {
			t.Error("NOT", expRes)
		}
	}
}

type TestFunc func(t *testing.T)

var i = 0

func runTest(t *testing.T, test TestFunc) {
	t.Run("test "+string(i), test)
	i++
}

func TestBasics(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testTrueToTrue", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testFalseToFalse", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testNotFalseToNotFalse", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testNotTrueToNotTrue", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testFalseToNotFalse", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testNotTrueToTrue", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testErrorTrueToNotTrue", "opposite conditions on K") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Basic/testNotFalseToTrue", "value of K is now True") })

	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/BasicChain/testNotFalseToTrueChain", "value of K is now True")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/BasicChain/testNotFalseToFalseChain", "value of K is now False")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/BasicChain/testErrorTrueToFalseChain", "opposite conditions on K")
	})
	runTest(t, func(t *testing.T) { testFile(t, "../tests/BasicChain/testTrueToTrueChain", "value of K is now True") })
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/BasicChain/testNotTrueToTrueChain", "value of K is now True")
	})
	runTest(t, func(t *testing.T) { testFile(t, "../tests/BasicChain/testFalseToTrueChain", "value of K is now True") })
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/BasicChain/testNotTrueToFalseChain", "opposite conditions on K")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/BasicChain/testFalseToFalseChain", "value of K is now False")
	})

	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/BasicChain/testNotFalseToTrueChainOtherOpposite", "opposite conditions on B")
	})
}

func TestRightOr(t *testing.T) {
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/RightSide/Or/testErrorTrueToNotOr", "opposite conditions on K")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/RightSide/Or/testOrOr", "value of J is now True")
	})
}

func TestRightNot(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/Not/testNot", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/Not/testNotNot", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/Not/testNotNot.1", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/Not/testNotNotNot", "value of K is now False") })
}

func TestRightAnd(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/And/testAnd", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/And/testAndFalse", "opposite conditions on K") })
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/RightSide/And/testAndOtherFalse", "opposite conditions on J")
	})
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/And/testAndAnd", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/And/testAndAnd.2", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/And/testNotAnd", "value of K is now Undetermined") })
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/RightSide/And/testNotAnd.1", "value of K is now False")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/RightSide/And/testNotAndAnd", "value of K is now Undetermined")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/RightSide/And/testNotAndAnd.1", "value of K is now False")
	})
}

func TestRightXor(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/Xor/testTrueXorVar.1", "value of J is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/Xor/testFalseXorVar", "value of J is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/RightSide/Xor/testTrueXorVar", "value of J is now False") })
}

func TestComplex(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.1", "value of E is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.2", "value of F is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.3", "value of G is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.4", "value of H is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.5", "value of I is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.6", "value of J is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.6bis", "value of J is now Undetermined") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.6.2", "opposite conditions on E") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.6.3", "value of J is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.7", "value of J is now Undetermined") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.8", "value of J is now Undetermined") }) // takes too long
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.9", "value of L is now Undetermined") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Complex/testComplex.9.1", "value of L is now Undetermined") })

}

func TestOptimizer(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/NotAndToOrNotNot", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/NotOrToAndNotNot", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/Xor", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/NotXor", "value of K is now False") })
}

func TestSpecial(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testSequential", "opposite conditions on B") })
}

func TestCorrection(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets", "value of E is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.1", "value of E is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.1.1", "value of E is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.2.1", "value of E is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.3.1", "value of E is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.4.1", "value of E is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.5.1", "value of E is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.6.1", "value of E is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.7.1", "value of E is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.8.1", "value of E is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Brackets.9.1", "value of E is now True") })

	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Not", "value of A is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Not.1", "value of A is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Not.1.1", "value of A is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Not.2.1", "value of A is now False") })

	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Or", "value of A is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Or.1", "value of A is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Or.1.1", "value of A is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Or.2.1", "value of A is now True") })

	runTest(t, func(t *testing.T) {
		testFileResults(t, "../tests/Correction/RightAnd", []string{
			"value of A is now True",
			"value of F is now True",
			"value of K is now True",
			"value of P is now True",
		})
	})
	runTest(t, func(t *testing.T) {
		testFileResults(t, "../tests/Correction/RightAnd.1", []string{
			"value of A is now True",
			"value of F is now True",
			"value of K is now False",
			"value of P is now True",
		})
	})

	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/SameMultiple", "value of A is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/SameMultiple.1", "value of A is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/SameMultiple.1.1", "value of A is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/SameMultiple.2.1", "value of A is now True") })

	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Xor", "value of A is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Xor.1", "value of A is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Xor.1.1", "value of A is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Correction/Xor.2.1", "value of A is now False") })

}

func TestParsing(t *testing.T) {
	folder := "../tests/Parsing/Condition/"
	runTest(t, func(t *testing.T) { testFile(t, folder+"testBadOperator", "Unknown char: -") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testRightBadOperator", "Unknown char: -") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testBadMiddle", "Bad syntax on line 1: No operator found") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testBadMiddle.1", "Bad syntax on line 1: No operator found") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testBadMiddle.1.1", "Unknown char: =") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testRightNothing", "Rule right side is empty!") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testLeftNothing", "Rule left side is empty!") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testLeftNothing.1", "Rule left side is empty!") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testOperatorFirst", "issue with +") })
	runTest(t, func(t *testing.T) { testFile(t, folder+"testMissingOperator", "Missing operator around A") })
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testMissingOperator.1", "! cannot be placed alone between operands")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testMissingOperator.1.1", "! cannot be placed alone between operands")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testMissingOperator.2.1", "! cannot be placed alone between operands")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testMissingOperand", "issue with +")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testMissingOperand.1", "+ operator requires two operands")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets", "missing operator between operands around A")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.1", "! cannot be placed alone between operands")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.1.1", "! cannot be placed alone between operands")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.2.1", "! cannot be placed alone between operands")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.3.1", "issue with +")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.4.1", "extra closing bracket")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.5.1", "extra closing bracket")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.6.1", "missing operator between operands around A")
	})
	runTest(t, func(t *testing.T) {
		testFile(t, folder+"testBadBrackets.7.1", "issue with +")
	})

}

func TestMainProgram(t *testing.T) {
	redirector.startRedirect()

	os.Args = []string{"testMain", "../tests/Correction/Xor..1"}
	runTest(t, func(t *testing.T) { main() })

	outStr := redirector.endRedirect()
	t.Log(outStr)
	expectedResult := "Le fichier spécifié est introuvable."
	if !strings.Contains(outStr, expectedResult) {
		t.Error("NOT", expectedResult)
	}
}
