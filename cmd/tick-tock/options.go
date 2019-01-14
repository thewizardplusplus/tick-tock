package main

import (
	"github.com/thewizardplusplus/tick-tock/interpreter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func parseOptions() interpreter.Options {
	kingpin.Version("v1.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('h')

	var options interpreter.Options
	kingpin.Flag("inbox", "Inbox buffer size.").Short('i').Default("10").IntVar(&options.InboxSize)
	kingpin.Flag("state", "Initial state.").
		Short('s').
		Default("__initialization__").
		StringVar(&options.InitialState)
	kingpin.Flag("message", "Initial message.").
		Short('m').
		Default("__initialize__").
		StringVar(&options.InitialMessage)
	kingpin.Arg("filename", `Source file name. Empty or "-" means stdin.`).
		StringVar(&options.Filename)
	kingpin.Parse()

	return options
}
