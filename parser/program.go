package parser

// Program ...
type Program struct {
	Definitions []*Definition `parser:"{ @@ }"`
}

// Definition ...
type Definition struct {
	Actor      *Actor      `parser:"@@"`
	ActorClass *ActorClass `parser:"| @@"`
}

// Actor ...
type Actor struct {
	Name       string   `parser:"\"actor\" @Ident"`
	Parameters []string `parser:"\"(\" [ @Ident { \",\" @Ident } [ \",\" ] ] \")\""`
	States     []*State `parser:"{ @@ } \";\""`
}

// ActorClass ...
type ActorClass struct {
	Name       string   `parser:"\"class\" @Ident"`
	Parameters []string `parser:"\"(\" [ @Ident { \",\" @Ident } [ \",\" ] ] \")\""`
	States     []*State `parser:"{ @@ } \";\""`
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
	Let        *LetCommand   `parser:"@@"`
	Start      *StartCommand `parser:"| @@"`
	Send       *SendCommand  `parser:"| @@"`
	Set        *SetCommand   `parser:"| @@"`
	Return     bool          `parser:"| @\"return\""`
	Expression *Expression   `parser:"| @@"`
}

// LetCommand ...
type LetCommand struct {
	Identifier string      `parser:"\"let\" @Ident \"=\""`
	Expression *Expression `parser:"@@"`
}

// StartCommand ...
type StartCommand struct {
	Name       *string         `parser:"\"start\" ( @Ident"`
	Expression *Expression     `parser:"| \"[\" @@ \"]\" )"`
	Arguments  ExpressionGroup `parser:"\"(\" @@ \")\""`
}

// SendCommand ...
type SendCommand struct {
	Name      string          `parser:"\"send\" @Ident"`
	Arguments ExpressionGroup `parser:"\"(\" @@ \")\""`
}

// SetCommand ...
type SetCommand struct {
	Name      string        `parser:"\"set\" @Ident"`
	Arguments []*Expression `parser:"\"(\" [ @@ { \",\" @@ } [ \",\" ] ] \")\""`
}
