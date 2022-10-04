# curate
Program for moving and renaming media files into folders

## Background

This is an attempt to make a more generic version of
[GardePro](https://github.com/madkins23/gardepro).
After I had that program I wanted to do something similar for
curating Google Photos files onto my NAS (which I do by hand, d'oh!).
Building a more generic version of GardePro seemed like an interesting project.

### Issues

The general issue is that I have a Network Attached Storage (NAS) device
onto which I curate images and videos that come from various sources.
I don't like to use the on-platform curation because it's always device-specific.
I prefer to arrange things into folders because that will survive moving to another device in the future.
So, yes, I'm making work for myself (D'oh!).

#### DCIM

_In my experience_ most devices that produce photos and videos store them
in a folder labelled DCIM, which stands for Digital Camera Images.
The contents of this folder are defined by the
[Design Rule for Camera File System](https://en.wikipedia.org/wiki/Design_rule_for_Camera_File_system)
(DCF) standard.

Within the DCIM folder are sub-folders which contain the actual media files.
When these sub-folders are deleted from the storage media (usually a memory card)
they are, _in my experience_, recreated by the device with the same names
and the files within are also recreated with the same names.
Copying the files to a new destination can then cause overwriting of
or conflicts with older files.

#### Deer Cameras

I have two GardePro deer cameras looking across the pasture behind my house.
These generate JPG and MP4 files in a fairly standard format
within the DCIM directory of an inserted memory card.

After I review the contents of a memory card I delete the DCIM sub-folders.
The cameras will then recreate the sub-folders the same names
and recreate contained media files with the same names.

While I only save a few files from each batch eventually I realized that
some new media files were overwriting older files with the same names.

For these cameras

* sub-folder names for these cameras are named with four digit numbers and
* file names are of the form `DSCF####.ext`

where

* `DSCF` likely stands for DSC File
* `####` is a number from `0001` to `9999`
* `ext` is the basic media type (e.g. `JPG` or `MP4`)

#### Google Photos

Google Photos are copied automatically into Google's cloud.
I generally download them from the website rather than trying to copy
them from my phone's DCIM directory (which is not on a memory chip).

_The following describes files from Google Pixel phones_.
Android phones by other manufacturers and/or telcoms
may use different naming conventions.

The files are currently named `PXL_YYYYMMDD_HHMMSSmmm[.type].ext` where:

* `PXL` stands for "Pixel" (the phone brand)
* `YYYYMMDD` represents the current date
* `HHMMSSmmm` represents the current time (`mmm` is milliseconds)
* `type` is an optional media type (e.g. `vr` or `PHOTOSPHERE`)
* `ext` is the basic media type (e.g. `JPG` or `MP4`)

The beginning bit is variable.
Google started with `IMG` and non-Google phones will have different strings.
I have also seen "type" information such as `PANO` in older files.

## Application Behavior

### Renaming

Source files are renamed to conform (more or less) to the way
Google names them `???_YYYYMMDD_HHMMSSmmm.ext`.
The `???` sequence varies depending on the source device.
The milliseconds in `mmm` will only be provided where convenient.

The milliseconds are currently not convenient,
so they are filled in with a three-digit string representing
an 8-bit CRC calculated from file contents.
This provides a constant number for a given file (as opposed to a random value)
and keeps files that just happen to be created in the same second from colliding.
For example, two deer cameras that happen to snap photos at the exact same second.

Source files are renamed based on the pattern of their original name.
If the pattern is not recognized the renaming fails and the file is not copied.
The `-unmatched` flag can be used to allow these files to be copied without renaming.

### Copying

When files are copied to the target directory the file names may collide.
In this case the file contents are compared to see if they are identical.
If they match then the copy is redundant and it is skipped.

If the file contents are _not_ identical then an error is flagged in the log
and via an alert box if the `-alert` flag is set.

### Early Termination on Error

If the `-alert` flag is set and a source file fails to copy for a reason
that shows the alert dialog the loop will terminate.
Without the `-alert` flag the loop will continue and all errors will be logged.

Having to click 'OK' on dozens of source files because of some
user error is _really_ annoying.

## Installation and Usage

This isn't the command-line usage which can be found in the
[application source](https://github.com/madkins23/curate/blob/main/cmd/curate/curate.go),
the [godoc](https://pkg.go.dev/github.com/madkins23/curate/cmd/curate),
or by building and running it without arguments or with the `-h` argument.
This section describes how I configure the application on my system.

When I thought about how I wanted to use this application,
I decided that the simplest thing would be to drag and drop
a file onto a desktop icon.
I created the application to work on a single source file at a time
and hooked it into a `.desktop` file within the `~/Desktop` directory.
This is my `~/Desktop/coyotes.desktop` file:

    [Desktop Entry]
    Version=0.1
    Type=Application
    Name=Curate
    Comment=Target for dropping photos and videos to the NAS
    Exec=/home/user/bin/curate -alert -target=/tmp/nasdir %F
    Terminal=false
    Categories=Utility;Application;

The file shows up on the desktop as a generic icon
since I didn't bother to configure a custom icon.
The file must have executable permission and
`Allow Launching` must be set in the right-click properties menu for the icon
(under Ubuntu, anyway, I have not tested this elsewhere).

When I drag and drop one or more files from the memory chip to the icon
the `curate` application is called with
the  file path(s) in place of the `%F` argument in the `Exec` string.
The application runs once for each drag and drop.

## Packages

Curate uses the following Go packages:

* [github.com/abema/go-mp4](https://github.com/abema/go-mp4) to get MP4 creation date/time
* [github.com/dsoprea/go-exif](https://github.com/dsoprea/go-exif) to get JPG creation date/time
* [github.com/madkins23/go-utils](https://github.com/madkins23/go-utils) for error messages and log file configuration
* [github.com/rs/zerolog](https://github.com/rs/zerolog) for pretty logging
* [github.com/sigurn/crc8](https://github.com/sigurn/crc8) to calculate three digit CRC numbers when milliseconds are not available
* [github.com/sqweek/dialog](https://github.com/sqweek/dialog)
  to display error messages directly to the user as they occur
  (they are also logged to a file)
* [github.com/udhos/equalfile](https://github.com/udhos/equalfile) to compare files
  in the case of duplicate target paths
* [github.com/sigurn/crc8](https://github.com/sigurn/crc8)
