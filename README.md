# ![](docs/logo/logo.png) Tick-tock

Interpreter of the Tick-tock programming language.

## Installation

```
$ go get github.com/thewizardplusplus/tick-tock/...
```

## Usage

```
$ tick-tock -v | --version
$ tick-tock -h | --help
$ tick-tock [-i SIZE | --inbox SIZE] [<filename>]
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

Copyright &copy; 2018 thewizardplusplus
