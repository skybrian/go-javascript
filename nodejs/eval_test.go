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
