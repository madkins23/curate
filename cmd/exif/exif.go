/*
exif is a tool for looking at EXIF values for a file.

This tool may evolve over time without regard to backwards compatibility
or proper version numbering.
Currently, the default behavior is to enumerate any EXIF properties in the file.
There is no other behavior at this time.

Usage:
    exif [flags]* [file]

Flags:
    -enumerate   show the EXIF properties in the file
*/
package main

import (
	"errors"
	"flag"
	"os"

	"github.com/madkins23/go-utils/log"

	EXIF "github.com/madkins23/curate/internal/exif"
)

func main() {
	log.Console()
	log.Info().Msg("exif starting")
	defer log.Info().Msg("exif finished")

	var enumerate bool

	flags := flag.NewFlagSet("Curate", flag.ContinueOnError)
	flags.BoolVar(&enumerate, "enumerate", true, "Show the EXIF properties in the file")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			log.Fatal().Err(err).Msg("Error parsing command line flags")
		} else {
			return
		}
	}

	if flags.NArg() < 1 {
		log.Fatal().Msg("Filename must be specified as last argument")
	}

	index, err := EXIF.GetIndex(flags.Arg(0))
	if err != nil {
		log.Fatal().Err(err).Msg("Unable to get EXIF index")
	}

	if enumerate {
		if err = EXIF.EnumerateIndex(index); err != nil {
			log.Fatal().Err(err).Msg("Unable to enumerate EXIF index")
		}
	}
}
