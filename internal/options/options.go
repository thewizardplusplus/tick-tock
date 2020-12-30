package options

import (
	"io"
	"path/filepath"
	"strconv"

	"github.com/thewizardplusplus/tick-tock/interpreter"
	"github.com/thewizardplusplus/tick-tock/runtime"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

// ...
const (
	Version = "v2.2.1"

	DefaultInboxSize      = 10
	DefaultInitialState   = "__initialization__"
	DefaultInitialMessage = "__initialize__"
)

// Dependencies ...
type Dependencies struct {
	UsageWriter io.Writer
	ErrorWriter io.Writer
	Exiter      runtime.Exiter
}

// Parse ...
func Parse(args []string, dependencies Dependencies) (interpreter.Options, error) {
	app := kingpin.New(filepath.Base(args[0]), "")
	app.UsageWriter(dependencies.UsageWriter)
	app.ErrorWriter(dependencies.ErrorWriter)
	app.Terminate(dependencies.Exiter)

	app.Version(Version)
	app.VersionFlag.Short('v')
	app.HelpFlag.Short('h')

	var options interpreter.Options
	app.Flag("inbox", "Inbox buffer size.").
		Short('i').
		Default(strconv.Itoa(DefaultInboxSize)).
		IntVar(&options.InboxSize)
	app.Flag("state", "Initial state.").
		Short('s').
		Default(DefaultInitialState).
		StringVar(&options.InitialState)
	app.Flag("message", "Initial message.").
		Short('m').
		Default(DefaultInitialMessage).
		StringVar(&options.InitialMessage)
	app.Arg("filename", `Source file name. Empty or "-" means stdin.`).StringVar(&options.Filename)

	if _, err := app.Parse(args[1:]); err != nil {
		return interpreter.Options{}, err
	}

	return options, nil
}
