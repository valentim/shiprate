package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"shiprate/internal/shipping"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("shiprate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	weight := fs.Float64("weight", 0, "weight in kg")
	length := fs.Float64("length", 0, "length in cm")
	width := fs.Float64("width", 0, "width in cm")
	height := fs.Float64("height", 0, "height in cm")
	zone := fs.String("zone", "", "destination zone, for example eu")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 1
	}
	if fs.NArg() > 0 {
		fmt.Fprintf(stderr, "shiprate: unexpected argument %q\n", fs.Arg(0))
		return 1
	}
	if missing := missingFlags(fs, "weight", "length", "width", "height", "zone"); len(missing) > 0 {
		fmt.Fprintf(stderr, "shiprate: missing required flags: %s\n", strings.Join(missing, ", "))
		return 1
	}

	quote, err := shipping.StandardRates().Quote(shipping.Parcel{
		WeightKg: *weight,
		LengthCm: *length,
		WidthCm:  *width,
		HeightCm: *height,
		Zone:     parseZone(*zone),
	})
	if err != nil {
		fmt.Fprintf(stderr, "shiprate: %v\n", err)
		return 1
	}

	if err := printText(stdout, quote); err != nil {
		fmt.Fprintf(stderr, "shiprate: %v\n", err)
		return 1
	}
	return 0
}

func parseZone(s string) shipping.Zone {
	return shipping.Zone(strings.ToLower(strings.TrimSpace(s)))
}

func missingFlags(fs *flag.FlagSet, required ...string) []string {
	set := map[string]bool{}
	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })

	var missing []string
	for _, name := range required {
		if !set[name] {
			missing = append(missing, "-"+name)
		}
	}
	return missing
}

func printText(w io.Writer, q shipping.Quote) error {
	for _, c := range q.Charges {
		if _, err := fmt.Fprintf(w, "%-15s %8s\n", c.Label, c.Amount); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintf(w, "%-15s %8s EUR\n", "total", q.Total)
	return err
}
