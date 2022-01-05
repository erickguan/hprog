/*
package main

import (
	"bufio"
	"fmt"
	"hprog/lexer"
	"hprog/token"
	"os"
)

func readline(idet string, scanner *bufio.Scanner) bool {
	fmt.Print(idet)
	return scanner.Scan()
}

func main() {
	fopen, _ := os.Open("./file.math")

	fbuffer := make([]byte, 1024)

	data, _ := fopen.Read(fbuffer)

	expression := string(fbuffer[:data])

	// TODO: get the tokens
	lex := lexer.Init(expression)
	for tkn := range lex.Consume() {
		// NOTE: i don't like this
		token.Print(tkn)
	}

	//st := stack.Stack{}

	defer fopen.Close()
}
*/

package main

import (
	"bufio"
	"flag"
	"fmt"
	"hprog/lexer"
	"hprog/token"
	"os"
	"strings"
)

func readline(idet string, scanner *bufio.Scanner) bool {
	fmt.Print(idet)
	return scanner.Scan()
}

func main() {
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

		fmt.Println("Hprog Version 0.0.0.0.0.0.0.0.2")
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
}
