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
	Let        *LetCommand   `parser:"@@"`
	Send       *SendCommand  `parser:"| @@"`
	Set        *string       `parser:"| \"set\" @Ident"`
	Out        *string       `parser:"| \"out\" ( @String | @RawString )"`
	Sleep      *SleepCommand `parser:"| @@"`
	Exit       bool          `parser:"| @\"exit\""`
	Expression *Expression   `parser:"| @@"`
}

// LetCommand ...
type LetCommand struct {
	Identifier string      `parser:"\"let\" @Ident \"=\""`
	Expression *Expression `parser:"@@"`
}

// SendCommand ...
type SendCommand struct {
	Name      string        `parser:"\"send\" @Ident"`
	Arguments []*Expression `parser:"\"(\" [ @@ { \",\" @@ } [ \",\" ] ] \")\""`
}

// SleepCommand ...
type SleepCommand struct {
	Minimum *float64 `parser:"\"sleep\" ( @Int | @Float )"`
	Maximum *float64 `parser:"\",\" ( @Int | @Float )"`
}
