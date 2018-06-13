package main

import (
	"io/ioutil"
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

type TestFunc func(t *testing.T)

var i = 0

func runTest(t *testing.T, test TestFunc) {
	t.Run("test "+string(i), test)
	i++
}

func TestBasics(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testTrueToTrue", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testFalseToFalse", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotFalseToNotFalse", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotTrueToNotTrue", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testFalseToNotFalse", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotTrueToTrue", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testErrorTrueToNotTrue", "opposite conditions on K") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotFalseToTrue", "value of K is now True") })

	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotFalseToTrueChain", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotFalseToFalseChain", "value of K is now False") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testErrorTrueToFalseChain", "opposite conditions on K") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testTrueToTrueChain", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotTrueToTrueChain", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testFalseToTrueChain", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testNotTrueToFalseChain", "opposite conditions on K") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testFalseToFalseChain", "value of K is now False") })

	runTest(t, func(t *testing.T) {
		testFile(t, "../tests/testNotFalseToTrueChainOtherOpposite", "opposite conditions on B")
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
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.1", "value of E is now False") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.2", "value of F is now True") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.3", "value of G is now True") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.4", "value of H is now False") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.5", "value of I is now True") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.6", "value of J is now True") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.6bis", "value of J is now Undetermined") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.6.2", "opposite conditions on E") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.6.3", "opposite conditions on E") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.7", "value of J is now Undetermined") })
	// runTest(t, func(t *testing.T) { testFile(t, "../tests/testComplex.8", "value of J is now Undetermined") }) // takes too long
}

func TestOptimizer(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/NotAndToOrNotNot", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/NotOrToAndNotNot", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/Xor", "value of K is now True") })
	runTest(t, func(t *testing.T) { testFile(t, "../tests/Optimizer/NotXor", "value of K is now False") })
}

func TestSpecial(t *testing.T) {
	runTest(t, func(t *testing.T) { testFile(t, "../tests/testSequential", "opposite conditins on B") })
}
