package parser

// ExpressionGroup ...
type ExpressionGroup struct {
	Expressions []*Expression `parser:"[ @@ { \",\" @@ } [ \",\" ] ]"`
}
