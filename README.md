# gofile

A simple utility to install go packages from a gofile (gofile.yml).

## Motivation

Ability to install system wide go packages.

A [similar project](https://github.com/Homebrew/homebrew-bundle) exists
for [Homebrew](https://brew.sh/).

## Por qué

Could time have been put to better use, by submitting these projects
into Homebrew?

> Probably, but that woudln't have been as much fun.

Could this have been a shell script?

> See above.

## Installation

    $  go get github.com/retr0h/gofile

## Usage

Create the gofile.

    $ cat gofile.yml
    ---
    - url: github.com/simeji/jid/cmd/jid
    - url: golang.org/x/lint/golint
    - url: golang.org/x/tools/cmd/goimports
    - url: github.com/arsham/figurine

Install go packages specified in the default gofile.yml.

    $ gofile install

Install go packages from an alternate gofile.

    $ gofile install --filename path/to/alternate.yml

[![asciicast](https://asciinema.org/a/192665.png)](https://asciinema.org/a/192665?speed=2&autoplay=1&loop=1)

## Dependencies

    $ go get github.com/golang/dep/cmd/dep

## Building

    $ make build
    $ tree .build/

## Testing

    $ make test

## License

MIT