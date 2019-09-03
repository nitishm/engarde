package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/nitishm/engarde/pkg/parser"
)

func main() {
	log.Printf("Reading input from STDIN. Use the pipe \"|\" operator to redirect traffic to engarde")
	parser := parser.New()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		accessLog, err := parser.Parse(line)
		if err != nil {
			accessLog.ParseError = err.Error()
		}
		bAccessLog, err := json.MarshalIndent(accessLog, "", "  ")
		if err != nil {
			accessLog.ParseError = err.Error()
		}
		fmt.Fprintf(os.Stdout, string(bAccessLog))
	}
}
