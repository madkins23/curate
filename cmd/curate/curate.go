/*
Curate moves (and optionally renames) media files into folders.

Usage:
  curate [flag]* [sourceFile]+

One or more [sourceFile] arguments can be provided.
Each [sourceFile] must be an absolute (non-relative) path.

The flags are:
  -alert
    	Show errors in alert panels in addition to logging them
  -console
    	Log to the console instead of the specified log file
  -debug uint
    	Debugging level (default 1)
  -logFile string
    	Log file path (default "/tmp/curate.log")
  -logJSON
    	Log output to file as JSON objects
  -normalize
    	Normalize basename(s) of source file(s) (default true)
  -target string
    	The destination directory for the media files
  -unmatched
    	When normalizing, copy files with unrecognized formats
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
	"github.com/madkins23/curate/internal/name"
)

var (
	alert     bool
	debug     uint
	normalize bool
	target    string
	unmatched bool
)

func main() {
	// Configure command-line flags.
	// TODO: Flag to separate files into years?
	flags := flag.NewFlagSet("Curate", flag.ContinueOnError)
	flags.BoolVar(&alert, "alert", false, "Show errors in alert panels in addition to logging them")
	flags.UintVar(&debug, "debug", 1, "Debugging level")
	flags.StringVar(&target, "target", "", "The destination directory for the media files")
	flags.BoolVar(&normalize, "normalize", true, "Normalize basename(s) of source file(s)")
	flags.BoolVar(&unmatched, "unmatched", false, "When normalizing, copy files with unrecognized formats")
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
	for _, source := range flags.Args() {
		log.Logger = baseLog.With().Str("source", source).Logger()
		if err := processSource(source); err != nil {
			if errors.Is(err, name.ErrNoPattern) {
				log.Warn().Str("reason", err.Error()).Msg("Skipping Copy")
			} else {
				log.Error().Err(err).Msg("Copy Failed")
			}
			if alert {
				dialog.Message(err.Error() + "\n" + source).Title("Process Error").Error()
				log.Info().Msg("Stop processing after first alert") // it's way too annoying.
				break
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

	var err error
	var basename string
	if normalize {
		// Normalize the source path:
		basename, err = name.Normalize(source)
		if err != nil {
			if !errors.Is(err, name.ErrNoPattern) {
				return fmt.Errorf("normalize base name: %w", err)
			} else if !unmatched {
				return err
			} else if debug > 0 {
				log.Info().Msg("Copying source with unmatched pattern")
			}
		}
	} else {
		basename = path.Base(source)
	}
	// Create path from target dir and basename.
	tgtPath := path.Join(target, basename)

	// Copy file:
	if err := ioutil.CopyFile(source, tgtPath); err != nil {
		if errors.Is(err, ioutil.ErrIdentical) {
			if debug > 0 {
				event := log.Info().Str("reason", err.Error())
				if normalize {
					event = event.Str("normalized", basename)
				}
				event.Msg("Skipping file copy")
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
