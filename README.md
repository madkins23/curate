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

* sub-folder names for these cameras are **TBD** and
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

**TBD**

This application renames an individual file and copies it to the NAS.
The renaming is done in a way that:

* separates the media into subdirectories by year,
* begins with the date and time so name ordering is chronological, and
* preserves the original basename (for no good reason).

I am counting on the two cameras not taking useful pictures
at the exact same second with the exact same base name.
There is code to check for non-identical overwrites.

### Copying

**TBD**

## Installation and Usage

**TBD**

This isn't the command-line usage which can be found in the
[application source](https://github.com/madkins23/curate/blob/main/cmd/curate/curate.go),
the [godoc](https://pkg.go.dev/github.com/madkins23/curate/cmd/curate),
or by building and running it without arguments.
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
    Exec=/home/marc/bin/curate -alert -target=/home/marc/tmp/nasdir %F
    Terminal=false
    Categories=Utility;Application;

The file shows up on the desktop as a generic icon
since I didn't bother to configure a custom icon.
The file must have executable permission and
`Allow Launching` must be set in the right-click properties menu.

When I drag and drop one or more files from the memory chip to the icon
the `curate` application is called with
the  file path(s) in place of the `%F` argument in the `Exec` string.
The application runs once for each drag and drop.

## Packages

Curate uses the following Go packages:

* [github.com/madkins23/go-utils/log](https://github.com/madkins23/go-utils) for log file configuration
* [github.com/rs/zerolog](https://github.com/rs/zerolog) for pretty logging
* [github.com/sqweek/dialog](https://github.com/sqweek/dialog)
  to display error messages directly to the user as they occur
  (they are also logged to a file)
* [github.com/udhos/equalfile](https://github.com/udhos/equalfile) to compare files
  in the case of duplicate target paths
