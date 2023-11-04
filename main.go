package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/flosch/pongo2/v6"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

var flagContextFile = flag.StringArrayP("contextfile", "f", nil, "YAML-file(s) to read context from")
var flagContext = flag.StringArrayP("context", "c", nil, "Additional context (key=value)")
var flagHelp = flag.BoolP("help", "h", false, "Show help")

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "%s [flags] <filename>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 || *flagHelp {
		usage()
	}

	ctx := pongo2.Context{}
	for _, f := range *flagContextFile {
		cf, err := os.ReadFile(f)
		if err != nil {
			log.Print("Could not read context file:", err)
			os.Exit(1)
		}
		err = yaml.Unmarshal(cf, &ctx)
		if err != nil {
			log.Print("Could not parse context file:", err)
			os.Exit(1)
		}
	}

	for _, c := range *flagContext {
		parts := strings.SplitN(c, "=", 2)
		if len(parts) != 2 {
			usage()
		}

		ctx[parts[0]] = parts[1]
	}

	ctx["filename"] = args[0]
	tmpl, err := pongo2.FromFile(args[0])
	if err != nil {
		log.Print("Could not load template:", err)
		os.Exit(1)
	}

	s, err := tmpl.Execute(ctx)
	if err != nil {
		log.Print("Could not load template:", err)
		os.Exit(1)
	}

	fmt.Print(s)
}
