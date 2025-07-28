package analysis_test

import (
	"testing"

	"vlogcc-lsp/analysis"
)

// Mock LoadFileContent to load from a fake in-memory FS.
var fakeFS = map[string]string{
	"/home/xsoder/project/main.code": `
require "lib/lib.code";

defun main(void) int {
	mut int x = 10;
	printf("Hello world");
}
`,
	"/home/xsoder/project/lib/lib.code": `
defun printf(void) int {
	// dummy printf
}
`,
}

func fakeLoadFileContent(path string) string {
	if content, ok := fakeFS[path]; ok {
		return content
	}
	return ""
}

func TestLoadImports(t *testing.T) {
	state := analysis.NewState()

	// Override LoadFileContent to use fake FS
	analysis.LoadFileContent = fakeLoadFileContent

	rootPath := "/home/xsoder/project/main.code"
	rootContent := fakeFS[rootPath]

	state.OpenDocument(rootPath, rootContent)

	// After OpenDocument, imported files should be loaded
	if _, ok := state.Documents["/home/xsoder/project/lib/lib.code"]; !ok {
		t.Fatalf("Expected imported file lib/lib.code to be loaded")
	}

	// The imported symbols should be in GlobalSymbols
	sym, ok := state.GlobalSymbols["printf"]
	if !ok {
		t.Fatalf("Expected symbol printf from imported file to be present in GlobalSymbols")
	}

	if sym.Type != "function" {
		t.Errorf("Expected printf symbol to be function, got %s", sym.Type)
	}
}

