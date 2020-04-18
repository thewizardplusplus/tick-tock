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
	Name     string     `parser:"\"message\" @Ident"`
	Commands []*Command `parser:"{ @@ } \";\""`
}

// Command ...
type Command struct {
	Let        *LetCommand `parser:"@@"`
	Send       *string     `parser:"| \"send\" @Ident"`
	Set        *string     `parser:"| \"set\" @Ident"`
	Expression *Expression `parser:"| @@"`
}

// LetCommand ...
type LetCommand struct {
	Identifier string      `parser:"\"let\" @Ident \"=\""`
	Expression *Expression `parser:"@@"`
}
