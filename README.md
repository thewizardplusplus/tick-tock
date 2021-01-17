# ![](docs/logo/logo.png) Tick-tock

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/tick-tock?status.svg)](https://godoc.org/github.com/thewizardplusplus/tick-tock)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/tick-tock)](https://goreportcard.com/report/github.com/thewizardplusplus/tick-tock)
[![Build Status](https://travis-ci.org/thewizardplusplus/tick-tock.svg?branch=master)](https://travis-ci.org/thewizardplusplus/tick-tock)
[![codecov](https://codecov.io/gh/thewizardplusplus/tick-tock/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/tick-tock)

Interpreter of the Tick-tock programming language.

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/tick-tock.git
$ cd tick-tock
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

Build the project:

```
$ go install ./...
```

## Usage

```
$ tick-tock -v | --version
$ tick-tock -h | --help
$ tick-tock [options] [<filename>]
```

Options:

- `-v`, `--version` &mdash; show application version;
- `-h`, `--help` &mdash; show application help;
- `-i SIZE`, `--inbox SIZE` &mdash; inbox buffer size (default: `10`);
- `-s STATE`, `--state STATE` &mdash; initial state (default: `__initialization__`);
- `-m MESSAGE`, `--message MESSAGE` &mdash; initial message (default: `__initialize__`).

Arguments:

- `<filename>` &mdash; source file name; empty or `-` means stdin.

## IDE support

- [Atom](http://atom.io/) plugin: [language-tick-tock](tools/atom-plugin/language-tick-tock).

## Docs

[Docs](docs/) (RU).

## License

The MIT License (MIT)

Copyright &copy; 2019-2021 thewizardplusplus
