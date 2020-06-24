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
	Name       string     `parser:"\"state\" @Ident"`
	Parameters []string   `parser:"\"(\" [ @Ident { \",\" @Ident } [ \",\" ] ] \")\""`
	Messages   []*Message `parser:"{ @@ } \";\""`
}

// Message ...
type Message struct {
	Name       string     `parser:"\"message\" @Ident"`
	Parameters []string   `parser:"\"(\" [ @Ident { \",\" @Ident } [ \",\" ] ] \")\""`
	Commands   []*Command `parser:"{ @@ } \";\""`
}

// Command ...
type Command struct {
	Let        *LetCommand  `parser:"@@"`
	Send       *SendCommand `parser:"| @@"`
	Set        *string      `parser:"| \"set\" @Ident"`
	Return     bool         `parser:"| @\"return\""`
	Expression *Expression  `parser:"| @@"`
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
