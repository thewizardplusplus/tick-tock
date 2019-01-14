package main

import kingpin "gopkg.in/alecthomas/kingpin.v2"

type options struct {
	filename  string
	inboxSize int
}

func parseOptions() options {
	kingpin.Version("v1.0")
	kingpin.CommandLine.VersionFlag.Short('v')
	kingpin.CommandLine.HelpFlag.Short('h')

	inboxSize := kingpin.Flag("inbox", "Inbox buffer size.").Short('i').Default("10").Int()
	filename := kingpin.Arg("filename", `Source file name. Empty or "-" means stdin.`).String()
	kingpin.Parse()

	return options{*filename, *inboxSize}
}
