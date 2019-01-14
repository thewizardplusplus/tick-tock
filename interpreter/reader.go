package interpreter

func isEmptyFilename(filename string) bool {
	return len(filename) == 0 || filename == "-"
}
