/*
Curate moves (and optionally renames) media files into folders.

Usage:

    curate [flag]* [sourceFile]+

One or more [sourceFile] arguments can be provided.
Each [sourceFile] must be an absolute (non-relative) path.

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

var (
	debug *uint
)

func main() {
	flags := flag.NewFlagSet("Curate", flag.ContinueOnError)
	debug = flags.Uint("debug", 1, "Debug level")

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

	baseLog := log.Logger

	// Handle the source files.
	for _, arg := range flags.Args() {
		log.Logger = baseLog.With().Str("source", arg).Logger()
		if err := processSource(arg); err != nil {
			log.Error().Err(err).Msg("Process Failed")
			dialog.Message(err.Error() + "\n" + arg).Title("Process Error").Error()
		}
	}
}

func processSource(source string) error {
	if *debug > 0 {
		log.Info().Msg("Processing")
	}
	return errors.New("Bugger!")
}
