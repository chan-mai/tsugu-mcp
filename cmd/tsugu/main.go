// Command tsugu はJSONから相続関係説明図PDFを生成するCLI
//
//	tsugu -in family.json -out chart.pdf [-era wareki|both|seireki]
//	tsugu < family.json > chart.pdf
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"tsugu-mcp/relationchart"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	in := flag.String("in", "", "入力JSONのパス (省略時は標準入力)")
	out := flag.String("out", "", "出力PDFのパス (省略時は標準出力)")
	era := flag.String("era", "wareki", "日付表記: wareki(和暦) | both(和暦+西暦) | seireki")
	flag.Parse()

	opt, err := optionsFromEra(*era)
	if err != nil {
		return err
	}

	data, err := readInput(*in)
	if err != nil {
		return err
	}

	pdf, err := relationchart.GenerateFromJSON(data, opt)
	if err != nil {
		return err
	}

	return writeOutput(*out, pdf)
}

func optionsFromEra(era string) (relationchart.Options, error) {
	switch era {
	case "wareki", "":
		return relationchart.Options{Era: relationchart.EraWareki}, nil
	case "both":
		return relationchart.Options{Era: relationchart.EraWarekiWithSeireki}, nil
	case "seireki":
		return relationchart.Options{Era: relationchart.EraSeireki}, nil
	default:
		return relationchart.Options{}, fmt.Errorf("unknown -era value: %q (wareki|both|seireki)", era)
	}
}

func readInput(path string) ([]byte, error) {
	if path == "" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(path)
}

func writeOutput(path string, pdf []byte) error {
	if path == "" {
		_, err := os.Stdout.Write(pdf)
		return err
	}
	return os.WriteFile(path, pdf, 0o644)
}
