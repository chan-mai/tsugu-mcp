// Command tsuguはJSONから相続書類のPDFを生成するCLI
//
//		tsugu chart -in family.json -out chart.pdf [-era wareki|both|seireki]   # 相続関係説明図
//		tsugu touki -in touki.json  -out touki.pdf [-era wareki|both|seireki]   # 相続登記申請書
//	 tsugu bunkatsu -in bunkatsu.json -out bunkatsu.pdf [-era wareki|both|seireki] # 遺産分割協議書
//		cat x.json | tsugu chart > out.pdf
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/chan-mai/tsugu-mcp/agreement"
	"github.com/chan-mai/tsugu-mcp/internal/buildinfo"
	"github.com/chan-mai/tsugu-mcp/registration"
	"github.com/chan-mai/tsugu-mcp/relationchart"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: tsugu <chart|touki|bunkatsu|certificate|version> -in <json> -out <pdf> [-era wareki|both|seireki]")
	}
	switch args[0] {
	case "version", "-version", "--version", "-v":
		fmt.Println("tsugu", buildinfo.String())
		return nil
	case "chart":
		return runChart(args[1:])
	case "touki":
		return runTouki(args[1:])
	case "bunkatsu":
		return runBunkatsu(args[1:])
	case "certificate":
		return runCertificate(args[1:])
	default:
		return fmt.Errorf("unknown subcommand: %q (chart|touki|bunkatsu|certificate)", args[0])
	}
}

func runChart(args []string) error {
	in, out, era, err := parseFlags("chart", args)
	if err != nil {
		return err
	}
	data, err := readInput(in)
	if err != nil {
		return err
	}
	pdf, err := relationchart.GenerateFromJSON(data, relationchart.Options{Era: era})
	if err != nil {
		return err
	}
	return writeOutput(out, pdf)
}

func runTouki(args []string) error {
	in, out, era, err := parseFlags("touki", args)
	if err != nil {
		return err
	}
	data, err := readInput(in)
	if err != nil {
		return err
	}
	pdf, err := registration.GenerateFromJSON(data, registration.Options{Era: era})
	if err != nil {
		return err
	}
	return writeOutput(out, pdf)
}

func runBunkatsu(args []string) error {
	in, out, era, err := parseFlags("bunkatsu", args)
	if err != nil {
		return err
	}
	data, err := readInput(in)
	if err != nil {
		return err
	}
	pdf, err := agreement.GenerateFromJSON(data, agreement.Options{Era: era})
	if err != nil {
		return err
	}
	return writeOutput(out, pdf)
}

func runCertificate(args []string) error {
	in, out, era, err := parseFlags("certificate", args)
	if err != nil {
		return err
	}
	data, err := readInput(in)
	if err != nil {
		return err
	}
	pdf, err := agreement.GenerateCertificateFromJSON(data, agreement.Options{Era: era})
	if err != nil {
		return err
	}
	return writeOutput(out, pdf)
}

func parseFlags(name string, args []string) (in, out string, era relationchart.EraStyle, err error) {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	inF := fs.String("in", "", "入力JSONのパス(省略時は標準入力)")
	outF := fs.String("out", "", "出力PDFのパス(省略時は標準出力)")
	eraF := fs.String("era", "wareki", "日付表記: wareki | both | seireki")
	if err = fs.Parse(args); err != nil {
		return
	}
	era, err = parseEra(*eraF)
	return *inF, *outF, era, err
}

func parseEra(s string) (relationchart.EraStyle, error) {
	switch s {
	case "wareki", "":
		return relationchart.EraWareki, nil
	case "both":
		return relationchart.EraWarekiWithSeireki, nil
	case "seireki":
		return relationchart.EraSeireki, nil
	default:
		return 0, fmt.Errorf("unknown -era value: %q (wareki|both|seireki)", s)
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
