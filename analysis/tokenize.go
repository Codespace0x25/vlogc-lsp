package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	functionPattern     = regexp.MustCompile(`^\s*defun\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(.*\)\s*([a-zA-Z_][a-zA-Z0-9_]*)\s*{`)
	functionCallPattern = regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(.*\)\s*;`)
	variablePattern     = regexp.MustCompile(`^\s*mut\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\s*\*?)?)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*=`)
	constPattern        = regexp.MustCompile(`^\s*([a-zA-Z_][a-zA-Z0-9_]*(?:\s*\*?)?)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*=`)
	macroPattern        = regexp.MustCompile(`^\s*macro\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(.*\)\s*{`)
	importPattern       = regexp.MustCompile(`^\s*require\s+"([^"]+)"\s*;?\s*$`)
)

func Tokenize(document string) []string {
	return strings.Split(document, "\n")
}

func stripFileScheme(path string) string {
	const prefix = "file://"
	if after, ok :=strings.CutPrefix(path, prefix); ok  {
		return after
	}
	return path
}

func ParseImports(content string) []string {
	var imports []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if matches := importPattern.FindStringSubmatch(line); matches != nil {
			fmt.Fprintf(os.Stderr, "Found import: %s\n", matches[1])
			imports = append(imports, matches[1])
		}
	}
	return imports
}

func normalizePath(baseDir, imp string) string {
	imp = filepath.Clean(imp)
	if !filepath.IsAbs(imp) {
		imp = filepath.Join(baseDir, imp)
	}
	return filepath.Clean(imp)
}

func cleanURI(uri string) string {
	return filepath.Clean(uri)
}

var LoadFileContent = func(path string) string {
	path = stripFileScheme(path)
	path = filepath.Clean(path)

	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data)
		}
	}

	// fallback 
	homePath := filepath.Join("/home/xsoder/programming/jsonlang", path)
	if _, err := os.Stat(homePath); err == nil {
		data, err := os.ReadFile(homePath)
		if err == nil {
			return string(data)
		}
	}

	stdPath := filepath.Join("/usr/include", path)
	if _, err := os.Stat(stdPath); err == nil {
		data, err := os.ReadFile(stdPath)
		if err == nil {
			return string(data)
		}
	}

	fmt.Fprintf(os.Stderr, "Failed to load file: %s or %s or %s\n", path, homePath, stdPath)
	return ""
}

