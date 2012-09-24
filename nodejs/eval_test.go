package nodejs

import (
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	checkRunSilent(t, "")
	checkRunSilent(t, "123")
	checkRunSilent(t, "2+2")
	checkRun(t, "123\n", "console.log(123)")
}

func TestRunWithArgs(t *testing.T) {
	checkRunSilent(t, "", "")
	checkRun(t, "[ '" + nodePathCache + "' ]\n", "console.log(process.argv)")
	checkRun(t, "[ '1', '2', '3' ]\n", "console.log(process.argv.slice(1))", "1", "2", "3")
}

func TestRunReportsSyntaxErrors(t *testing.T) {
	checkRunSyntaxError(t, "Unexpected end of input", "[")
}

func TestEvalExpr(t* testing.T) {
	checkEval(t, "123", "123")
	checkEval(t, "5", "2+3")
	checkEval(t, `'hello'`, `"hello"`)
	checkEval(t, `[ 1, 2, 3 ]`, `[1,2,3]`)
	checkEval(t, `{ a: 1 }`, `{"a": 1}`)
}

func TestEvalExprReportsSyntaxErrors(t *testing.T) {
	checkEvalSyntaxError(t, "Unexpected token )", "[")
}

func TestEvalEach(t* testing.T) {
	checkEvalEach(t)
	checkEvalEach(t, "2+2", "4")
	checkEvalEach(t, "2+2", "4", "2+3", "5")
	checkEvalEach(t, `'hello'`, `'hello'`, `[1, 2, 3]`, `[ 1, 2, 3 ]`, `{"a": 1}`, `{ a: 1 }`)
}

// === end of tests ===

func checkRunSilent(t *testing.T, script string, args ...string) {
	checkRun(t, "", script, args...)
}

func checkRun(t *testing.T, expectedOut, script string, args ...string) {
	out, err := Run(script, args...)
	checkSuccess(t, expectedOut, script, out, err)
}

func checkEval(t *testing.T, expectedOut, expr string) {
	out, err := EvalExpr(expr)
	checkSuccess(t, expectedOut, expr, out, err)
}

func checkSuccess(t *testing.T, expectedOut, in, out string, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v => %v", in, err)
		return
	}
	if out != expectedOut {
		t.Errorf("unexpected output: %v => %v", in, out)
	}
}

func checkRunSyntaxError(t *testing.T, expectedError string, script string, args ...string) {
	out, err := Run(script, args...)
	checkSyntaxError(t, expectedError, script, out, err)
}

func checkEvalSyntaxError(t *testing.T, expectedErr, expr string) {
	out, err := EvalExpr(expr)
	checkSyntaxError(t, expectedErr, expr, out, err)
}

func checkSyntaxError(t *testing.T, expectedError, in, out string, err error) {
	if err == nil {
		t.Errorf("expected an error, got output: %v => %v", in, out)
	}
	actual := err.Error()
	if !strings.Contains(actual, "\nSyntaxError: " + expectedError + "\n") {
		t.Errorf("unexpected error: %v => %v", in, actual)
	}
}

func checkEvalEach(t *testing.T, args... string) {
	if len(args) % 2 != 0 {
		t.Fatalf("wrong number of args to checkEvalEach: %v", len(args))
	}
	ins := []string {}
	expected := []string {}
	for i := 0; i < len(args); i += 2 {
		ins = append(ins, args[i])
		expected = append(expected, args[i + 1])
	}
	outs, err := EvalEach(ins...)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	checkSlice(t, ins, outs, expected)
}

func checkSlice(t *testing.T, ins, outs, expectedOuts []string) {
	if len(ins) != len(expectedOuts) {
		t.Fatalf("number of expected outputs (%v) doesn't match number of inputs: %v", len(ins), len(expectedOuts))
	} else if len(expectedOuts) != len(outs) {
		t.Errorf("expected %v outputs; got %v", len(expectedOuts), len(outs))
		return
	}
	for i := range outs {
		if outs[i] != expectedOuts[i] {
			t.Errorf("unexpected output: %v -> %v", ins[i], outs[i])
		}
	}
}
