/*
Curate moves (and optionally renames) media files into folders.

Usage:

    curate [flags]

The flags are:

    -console
        Log to the console instead of the specified log file [false]
    -logFile
        Log file path [/tmp/curate.log]
    -logJSON
        Log output as JSON items [false]
    -debug
        Debugging level [1] (0..1)
*/
package main

import (
	"errors"
	"flag"
	"os"

	utilog "github.com/madkins23/go-utils/log"
	"github.com/rs/zerolog/log"
	"github.com/sqweek/dialog"
)

func main() {
	flags := flag.NewFlagSet("gardepro", flag.ContinueOnError)
	debug := flags.Uint("debug", 1, "Debug level")

	cof := utilog.ConsoleOrFile{}
	cof.AddFlagsToSet(flags, "/tmp/curate.log")

	if err := flags.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			dialog.Message(err.Error()).Title("Error parsing command line flags").Error()
		}
		return
	}

	if err := cof.Setup(); err != nil {
		dialog.Message(err.Error()).Title("Log File Creation").Error()
		return
	}
	defer cof.CloseForDefer()

	if *debug > 0 {
		log.Info().Msg("Curate starting")
		defer log.Info().Msg("Curate finished")
	}
}
