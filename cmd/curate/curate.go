/*
Curate moves (and optionally renames) media files into folders.

Usage:

    curate [flag]* [sourceFile]+

One or more [sourceFile] arguments can be provided.
Each [sourceFile] must be an absolute (non-relative) path.

The flags are:

    -target=[directory path]
        The destination directory for the media files (required)
    -alert
        Show errors in alert panels in addition to logging them [false]
    -console
        Log to the console instead of the specified log file [false]
    -alert
        Show error
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
	"fmt"
	"os"
	"path"

	utilLog "github.com/madkins23/go-utils/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sqweek/dialog"

	"github.com/madkins23/curate/internal/ioutil"
)

var (
	alert  bool
	debug  uint
	target string
)

func main() {
	// Configure command-line flags.
	flags := flag.NewFlagSet("Curate", flag.ContinueOnError)
	flags.BoolVar(&alert, "alert", false, "Show error alerts")
	flags.UintVar(&debug, "debug", 1, "Debug level")
	flags.StringVar(&target, "target", "", "Destination directory")
	cof := utilLog.ConsoleOrFile{}
	cof.AddFlagsToSet(flags, "/tmp/curate.log")
	if err := flags.Parse(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			if alert {
				dialog.Message(err.Error()).Title("Error parsing command line flags").Error()
			} else {
				fmt.Printf("Error parsing command line flags: %s\n", err)
			}
		}
		return
	}

	// Configure logging.
	if err := cof.Setup(); err != nil {
		if alert {
			dialog.Message(err.Error()).Title("Log File Creation").Error()
		} else {
			fmt.Printf("Log file creation error: %s\n", err)
		}
		return
	}
	defer cof.CloseForDefer()

	if debug > 0 {
		log.Info().Msg("Curate starting")
		defer log.Info().Msg("Curate finished")
	}

	if target == "" {
		fatalError("No target path", nil, nil)
	} else if err := ioutil.CheckDir(target); err != nil {
		fatalError("Bad target path", err, func(event *zerolog.Event) *zerolog.Event {
			return event.Str("target", target)
		})
	}

	// Create a log with the target string added to each log entry.
	baseLog := log.Logger.With().Str("target", target).Logger()

	// Handle the source files.
	for _, arg := range flags.Args() {
		log.Logger = baseLog.With().Str("source", arg).Logger()
		if err := processSource(arg); err != nil {
			log.Error().Err(err).Msg("Process Failed")
			if alert {
				dialog.Message(err.Error() + "\n" + arg).Title("Process Error").Error()
			}
		}
	}
}

// fatalError logs the specified message and any non-nil error,
// optionally displaying an alert window to the user in real-time,
// then exits the program.
func fatalError(message string, err error, extra func(*zerolog.Event) *zerolog.Event) {
	msg := message
	if err != nil {
		msg += ":\n" + err.Error()
	}

	if alert {
		dialog.Message(msg).Title("Fatal Error").Error()
	}

	// Fatal() will call os.Exit() after logging, skipping defer statements in main().
	event := log.Fatal()
	if err != nil {
		event = event.Err(err)
	}
	if extra != nil {
		event = extra(event)
	}
	event.Msg(message)
}

// processSource optionally renames and then copies a single source file to the target directory.
func processSource(source string) error {
	if debug > 0 {
		log.Info().Msg("Processing")
	}

	// TODO: Normalize (uniqify?) name if necessary.
	//  This may include creating a subdirectory.
	//  Does this include verifying/creating the subdirectory
	//   or is that done later?

	// Create path from target dir and (subdirectory/)name.
	tgtPath := path.Join(target, path.Base(source))

	// Copy file:
	if err := ioutil.CopyFile(source, tgtPath); err != nil {
		if errors.Is(err, ioutil.ErrIdentical) {
			if debug > 0 {
				log.Info().Str("reason", err.Error()).Msg("Skipping file copy")
			}
		} else {
			fatalError("Copy error", err, func(event *zerolog.Event) *zerolog.Event {
				return event.Str("target-path", tgtPath)
			})
		}
	} else if debug > 0 {
		log.Info().Str("target", tgtPath).Msg("Copied file")
	}

	// TODO: Delete file after copy?

	return nil
}

//var errUnknownExtension = errors.New("unknown file extension")
//
//switch ext := strings.ToLower(filepath.Ext(source)); ext {
//case ".jpg", ".jpeg":
//	log.Debug().Msg("TBD JPG")
//case ".mp4":
//	log.Debug().Msg("TBD MP4")
//default:
//	return errUnknownExtension
//}
