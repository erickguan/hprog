package vm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func BenchmarkVM(b *testing.B) {
	var testCases = [...]string{
		"./data_bench/bool_op.hp",
		"./data_bench/print_op.hp",
		"./data_bench/sub_ops.hp",
		"./data_bench/string_ops.hp",
		"./data_bench/decl_op.hp",
	}
	for _, inputFile := range testCases {
		var buffer []string

		hFile, err := os.Open(inputFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		fileScanner := bufio.NewScanner(hFile)
		for fileScanner.Scan() {
			buffer = append(buffer, fileScanner.Text())
		}

		v := VM{}
		v.InitVM()
		status := v.Interpret(strings.Join(buffer[:], "\n"))
		if status != INTER_OK {
			fmt.Println("Runtime error.")
			break
		}
	}
}

func _TestNumbers(t *testing.T) {
	var testCases = []string{
		"1 + \"3\"",
		"\"1\" + 3",
	}
	for _, v := range testCases {
		Execute(v, t)
	}
}

func Execute(expression string, t *testing.T) {
	v := VM{}
	v.InitVM()
	status := v.Interpret(expression)
	if status != INTER_OK {
		t.Errorf("input %s", expression)
	}
}
