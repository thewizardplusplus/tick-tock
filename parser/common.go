package parser

// IdentifierGroup ...
type IdentifierGroup struct {
	Identifiers []string `parser:"[ @Ident { \",\" @Ident } [ \",\" ] ]"`
}

// ExpressionGroup ...
type ExpressionGroup struct {
	Expressions []*Expression `parser:"[ @@ { \",\" @@ } [ \",\" ] ]"`
}
