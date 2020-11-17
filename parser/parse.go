package parser

// Parse ...
func Parse(code string) (*Program, error) {
	program := new(Program)
	if err := ParseToAST(code, program); err != nil {
		return nil, err
	}

	return program, nil
}
