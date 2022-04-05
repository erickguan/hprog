package main

import (
	"bufio"
	"fmt"

	"github.com/badc0re/hprog/lexer"
	"github.com/badc0re/hprog/token"
)

func readline(idet string, scanner *bufio.Scanner) bool {
	fmt.Print(idet)
	return scanner.Scan()
}

func main() {
	/*
		// direct file parser
		//test()
		var buffer []string
		var inputFile string
		flag.StringVar(&inputFile, "file", "", "Input hell file.")
		flag.Parse()

		fmt.Println(inputFile)
		hFile, err := os.Open(inputFile)

		if hFile != nil {
			if err != nil {
				fmt.Println(err)
				return
			}
			fileScanner := bufio.NewScanner(hFile)
			for fileScanner.Scan() {
				buffer = append(buffer, fileScanner.Text())
			}

			lex := lexer.Init(strings.Join(buffer[:], "\n"))
			for tkn := range lex.Consume() {
				token.Print(tkn)
			}
		} else {
			const idet = "hprog> "

			fmt.Println("Hprog Version 0.01")
			fmt.Println("One way to escape, ctr-c to exit.")

			scanner := bufio.NewScanner(os.Stdin)

			onNewLine := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
				return bufio.ScanLines(data, atEOF)
			}

			scanner.Split(onNewLine)
			for {
				for readline(idet, scanner) {
					var sline = scanner.Text()
					if len(sline) > 0 {
						var tokens []token.Token
						lex := lexer.Init(sline)
						for tkn := range lex.Consume() {
							tokens = append(tokens, tkn)
							if tkn.Type == token.EOF || tkn.Type == token.ERR {
								break
							}
							token.Print(tkn)
						}

						buffer = append(buffer, sline)
					}
				}
			}

			if scanner.Err() != nil {
				fmt.Printf("error: %s\n", scanner.Err())
			}
		}
	*/
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
			status := v.Interpret(".11")
			if status == vm.INTER_RUNTIME_ERROR {
				fmt.Println("Runtime error.")
			}
		v.FreeVM()
	*/
	lex := lexer.Init("(a)")

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
}
