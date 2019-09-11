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
	limit            int
)

func init() {
	flag.BoolVar(&isIstioProxyLogs, "use-istio", false, "Enable to use istio-proxy format")
	flag.IntVar(&limit, "limit", 0, "Limit number of lines parsed. Set to 0 for unlimited scroll")
}

func main() {
	var p *parser.Parser
	var enableLimits bool

	flag.Parse()

	log.Printf("Reading input from STDIN. Use the pipe \"|\" operator to redirect traffic to engarde")
	if isIstioProxyLogs {
		p = parser.New(parser.IstioProxyAccessLogsPattern)
	} else {
		p = parser.New(parser.EnvoyAccessLogsPattern)
	}

	if limit > 0 {
		enableLimits = true
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

		if enableLimits {
			limit = limit - 1
			if limit == 0 {
				os.Exit(0)
			}
		}
	}
}
