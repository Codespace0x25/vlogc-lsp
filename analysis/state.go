package analysis

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"vlogcc-lsp/lsp"
)

type State struct {
	Documents     map[string]string
	Symbols       map[string]lsp.Symbol
	GlobalSymbols map[string]lsp.Symbol
}

func NewState() State {
	return State{
		Documents:     map[string]string{},
		Symbols:       map[string]lsp.Symbol{},
		GlobalSymbols: map[string]lsp.Symbol{},
	}
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	symbol, err := s.GetSymbolAtPosition(uri, position)
	if err != nil {
		return lsp.HoverResponse{
			Response: lsp.Response{
				RPC: "2.0",
				ID:  &id,
			},
			Result: lsp.HoverResult{
				Contents: "Error retrieving symbol",
			},
		}
	}

	contents := fmt.Sprintf(
		"Symbol: %s\nKind: %s\nType: %s\nMutable: %t",
		symbol.Name,
		symbol.Type,
		symbol.DataType,
		symbol.Mutable,
	)

	if symbol.Parameters != "" {
		contents += fmt.Sprintf("\nParameters: %s", symbol.Parameters)
	}

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: contents,
		},
	}
}

func (s *State) GetSymbolAtPosition(uri string, position lsp.Position) (lsp.Symbol, error) {
	document := s.Documents[uri]
	lines := Tokenize(document)

	if position.Line < 0 || position.Line >= len(lines) {
		return lsp.Symbol{}, fmt.Errorf("invalid position")
	}

	line := lines[position.Line]
	if line == "" {
		return lsp.Symbol{}, fmt.Errorf("empty line")
	}
	fmt.Fprintf(os.Stderr, "Processing line: %d -> %s\n", position.Line, line)

	if matches := functionPattern.FindStringSubmatch(line); matches != nil {
		return lsp.Symbol{Name: matches[1], Type: "function", DataType: matches[2]}, nil
	}
	if matches := variablePattern.FindStringSubmatch(line); matches != nil {
		return lsp.Symbol{Name: matches[2], Type: "variable", DataType: matches[1], Mutable: true}, nil
	}
	if matches := constPattern.FindStringSubmatch(line); matches != nil {
		return lsp.Symbol{Name: matches[2], Type: "variable", DataType: matches[1], Mutable: false}, nil
	}
	if matches := functionCallPattern.FindStringSubmatch(line); matches != nil {
		// Check in global symbols
		if sym, ok := s.GlobalSymbols[matches[1]]; ok {
			return sym, nil
		}
		return lsp.Symbol{Name: matches[1], Type: "call", DataType: "unknown"}, nil
	}

	return lsp.Symbol{}, fmt.Errorf("no symbol found at the specified position")
}

func ParseSymbols(source string) []lsp.Symbol {
	var symbols []lsp.Symbol
	lines := strings.Split(source, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if matches := functionPattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, lsp.Symbol{
				Name:     matches[1],
				Type:     "function",
				DataType: matches[2],
				Mutable:  false,
			})
			continue
		}

		if matches := macroPattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, lsp.Symbol{
				Name:    matches[1],
				Type:    "macro",
				Mutable: false,
			})
			continue
		}

		if matches := importPattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, lsp.Symbol{
				Name:     matches[1],
				Type:     "import",
				DataType: "std-library",
				Mutable:  false,
			})
			continue
		}

		if matches := variablePattern.FindStringSubmatch(line); matches != nil {
			symbols = append(symbols, lsp.Symbol{
				Name:     matches[2],
				Type:     "variable",
				DataType: matches[1],
				Mutable:  true,
			})
			continue
		}

		if matches := constPattern.FindStringSubmatch(line); matches != nil && !strings.HasPrefix(line, "mut ") {
			symbols = append(symbols, lsp.Symbol{
				Name:     matches[2],
				Type:     "variable",
				DataType: matches[1],
				Mutable:  false,
			})
			continue
		}
	}

	return symbols
}

func (s *State) OpenDocument(uri, text string) {
	uri = stripFileScheme(uri)
	uri = cleanURI(uri)
	s.Documents[uri] = text
	if err := s.LoadImports(uri); err != nil {
		fmt.Fprintf(os.Stderr, "LoadImports error: %v\n", err)
	}
}

func (s *State) UpdateDocument(uri, text string) {
	uri = stripFileScheme(uri)
	uri = cleanURI(uri)
	s.Documents[uri] = text
	if err := s.LoadImports(uri); err != nil {
		fmt.Fprintf(os.Stderr, "LoadImports error: %v\n", err)
	}
}

func (s *State) LoadImports(uri string) error {
	uri = stripFileScheme(uri)
	uri = filepath.Clean(uri)

	content, ok := s.Documents[uri]
	if !ok {
		return fmt.Errorf("document not found: %s", uri)
	}

	imports := ParseImports(content)
	baseDir := filepath.Dir(uri)

	for _, imp := range imports {
		imp = filepath.Clean(imp)

		var impPath string
		if filepath.IsAbs(imp) {
			impPath = imp
		} else {
			impPath = filepath.Join(baseDir, imp)
		}
		impPath = filepath.Clean(impPath)

		if _, loaded := s.Documents[impPath]; !loaded {
			impContent := LoadFileContent(impPath)
			if impContent == "" {
				fmt.Fprintf(os.Stderr, "Failed to load import file content: %s\n", impPath)
				continue
			}
			s.Documents[impPath] = impContent

			symbols := ParseSymbols(impContent)
			for _, sym := range symbols {
				s.GlobalSymbols[sym.Name] = sym
			}

			if err := s.LoadImports(impPath); err != nil {
				fmt.Fprintf(os.Stderr, "LoadImports error for %s: %v\n", impPath, err)
			}
		}
	}
	return nil
}
