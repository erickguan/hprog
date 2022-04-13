package vm

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func BenchmarkLexerNumbers(b *testing.B) {
	var testCases = [...]string{
		"./data/bool_op.hp",
		"./data/print_op.hp",
		"./data/sub_ops.hp",
		"./data/string_ops.hp",
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
