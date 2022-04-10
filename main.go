package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/badc0re/hprog/vm"
)

func readline(idet string, scanner *bufio.Scanner) bool {
	//	fmt.Print(idet)
	return scanner.Scan()
}

func loadFile(inputFile string) {
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

	v := vm.VM{}
	v.InitVM()
	status := v.Interpret(strings.Join(buffer[:], "\n\n"))
	if status == vm.INTER_RUNTIME_ERROR {
		fmt.Println("Runtime error.")
	}
}

func OnNewLine(data []byte, atEOF bool) (advance int, token []byte, err error) {
	return bufio.ScanLines(data, atEOF)
}

func main() {
	var buffer []string
	var inputFile string

	flag.StringVar(&inputFile, "file", "", "Input hprog file.")
	flag.Parse()

	if len(inputFile) != 0 {
		loadFile(inputFile)
		os.Exit(1)
	}

	const idet = "hprog> "
	fmt.Println("Hprog Version 0.01")
	fmt.Println("One way to escape, ctr-c to exit.")

	// INPUT SCANNER
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(OnNewLine)

	// INIT VM
	v := vm.VM{}
	v.InitVM()

	// readlines and process
	for readline(idet, scanner) {
		var sline = scanner.Text()

		if scanner.Err() != nil {
			fmt.Printf("error: %s\n", scanner.Err())
		}

		if len(sline) > 0 {
			status := v.Interpret(sline)
			if status == vm.INTER_RUNTIME_ERROR {
				fmt.Println("Runtime error.")
			}
			buffer = append(buffer, sline)
		}
	}

	/*
		v.InitVM()
		chk := chunk.Chunk{}

		id := chk.AddVariable(123)
		chk.WriteChunk(codes.INSTRUC_CONSTANT, 1)
		chk.WriteChunk(id, 1)

		id2 := chk.AddVariable(456)
		chk.WriteChunk(codes.INSTRUC_CONSTANT, 1)
		chk.WriteChunk(id2, 1)

		chk.WriteChunk(codes.INSTRUC_ADDITION, 1)

		//chk.WriteChunk(codes.OP_NEGATE, 1)

		id3 := chk.AddVariable(1000)
		chk.WriteChunk(codes.INSTRUC_CONSTANT, 1)
		chk.WriteChunk(id3, 1)

		chk.WriteChunk(codes.INSTRUC_ADDITION, 1)

		chk.WriteChunk(codes.INSTRUC_RETURN, 1)

		chunk.DissasChunk(&chk, "simple instruction")
	*/
	// freeChunk(&chk)
	// freeChunk(&chk)

	/*
		v := vm.VM{}
		v.InitVM()
		status := v.Interpret("print1")
		if status == vm.INTER_RUNTIME_ERROR {
			fmt.Println("Runtime error.")
		}
		v.FreeVM()
			lex := lexer.Init("print1\n")
			var result []token.Token
			for {
				tkn, _ := lex.Consume()
				a := *tkn
				if a.Type == token.EOF {
					fmt.Println("DONE scan")
					break
				} else if a.Type == token.ERR {
					fmt.Println("ERROR scan")
					break
				}
				result = append(result, a)
			}
			for _, i := range result {
				fmt.Println("END RESULT:", token.ReversedTokenMap[i.Type], "value:", i.Value)
			}
	*/
}
