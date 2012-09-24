// Package nodejs runs JavaScript by invoking Node.js as a separate process.
// The "node" command must be available in the current process's path.
package nodejs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

var nodePathCache string

func findNode() (string, error) {
	if nodePathCache == "" {
		path, err := exec.LookPath("node")
		if err != nil {
			return "", err
		}
		nodePathCache = path
	}
	return nodePathCache, nil
}

// Run invokes Node with a script and waits until it exits.
// If Node exits successfully, it returns stdout and discards stderr.
// The script is passed in using -e, so command-line limits apply.
func Run(script string, args ...string) (string, error) {
	nodePath, err := findNode()
	if err != nil {
		return "", err
	}
	args = append([]string {"-e", script}, args...)
	cmd := exec.Command(nodePath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	switch e := err.(type) {
	case nil:
		return string(stdout.Bytes()), nil
	case *exec.ExitError:
	  	return "", &ExitError{e.ProcessState, &stderr};
	}
	return "", err
}

// EvalExpr returns the value of a JavaScript expression.
// The value will be formatted according to Node's util.inspect.
func EvalExpr(expr string) (string, error) {
	out, err := Run(`console.log(require("util").inspect(eval(process.argv[1])))`, "(" + expr + ")")
  	if err != nil {
		return "", err
	}
	return out[:len(out) - 1], nil
}

// EvalEach returns the value of each JavaScript expression, formatted according to Node's util.inspect.
func EvalEach(exprs... string) ([]string, error) {
	script := `
var util = require("util");
var results = [];
for (var i = 1; i < process.argv.length; i++) {
	results[i-1] = util.inspect(eval(process.argv[i]));
}
console.log(JSON.stringify(results));
`
	args := []string {}
	for _, expr := range exprs {
		args = append(args, "(" + expr + ")")
	}

	out, err := Run(script, args...)
  	if err != nil {
		return nil, err
	}
	outs := []string {}
	err = json.Unmarshal([]byte(out), &outs)
	if err != nil {
		return nil, err
	}
	return outs, nil
}

type ExitError struct {
	*os.ProcessState
	Stderr *bytes.Buffer
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("node.js script failed:\n%v", e.Stderr)
}
