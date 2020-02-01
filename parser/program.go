package parser

// Program ...
type Program struct {
	Actors []*Actor `parser:"{ @@ }"`
}

// Actor ...
type Actor struct {
	States []*State `parser:"\"actor\" { @@ } \";\""`
}

// State ...
type State struct {
	Name     string     `parser:"\"state\" @Ident"`
	Messages []*Message `parser:"{ @@ } \";\""`
}

// Message ...
type Message struct {
	Name       string     `parser:"\"message\" @Ident"`
	Parameters []string   `parser:"\"(\" [ @Ident { \",\" @Ident } [ \",\" ] ] \")\""`
	Commands   []*Command `parser:"{ @@ } \";\""`
}

// Command ...
type Command struct {
	Send  *string       `parser:"\"send\" @Ident"`
	Set   *string       `parser:"| \"set\" @Ident"`
	Out   *string       `parser:"| \"out\" ( @String | @RawString )"`
	Sleep *SleepCommand `parser:"| @@"`
	Exit  bool          `parser:"| @\"exit\""`
}

// SleepCommand ...
type SleepCommand struct {
	Minimum *float64 `parser:"\"sleep\" ( @Int | @Float )"`
	Maximum *float64 `parser:"\",\" ( @Int | @Float )"`
}
