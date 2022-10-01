# curate
Program for moving and renaming media files into folders

## Background

This is an attempt to make a more generic version of
[GardePro](https://github.com/madkins23/gardepro).
After I had that program I wanted to do something similar for
curating Google Photos files onto my NAS (which I do by hand, d'oh!).
Building a more generic version of GardePro seemed like an interesting project.

## Packages

Curate uses the following Go packages:

* [github.com/madkins23/go-utils/log](https://github.com/madkins23/go-utils) for log file configuration
* [github.com/rs/zerolog](https://github.com/rs/zerolog) for pretty logging
* [github.com/sqweek/dialog](https://github.com/sqweek/dialog)
  to display error messages directly to the user as they occur
  (they are also logged to a file)
* [github.com/udhos/equalfile](https://github.com/udhos/equalfile) to compare files
  in the case of duplicate target paths
