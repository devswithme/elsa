package generate

// FindGoModDir searches upward from a starting path to find the directory containing go.mod
// This function is useful for locating the Go module root from any file within the module
// It traverses up the directory tree until it finds a go.mod file or reaches the filesystem root
func (g *Generator) FindGoModDir(start string) (string, error) {
	return findGoModDir(start)
}
