package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nitishm/engarde/pkg/parser"
)

var (
	isIstioProxyLogs bool
)

func init() {
	flag.BoolVar(&isIstioProxyLogs, "use-istio", false, "Enable to use istio-proxy format")
}

func main() {
	var p *parser.Parser
	flag.Parse()

	log.Printf("Reading input from STDIN. Use the pipe \"|\" operator to redirect traffic to engarde")
	if isIstioProxyLogs {
		p = parser.New(parser.IstioProxyAccessLogsPattern)
	} else {
		p = parser.New(parser.EnvoyAccessLogsPattern)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		accessLog, err := p.Parse(line)
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
