package analysis

import (
	"fmt"
	"vlogcc-lsp/lsp"
)

type State struct {
	Documents map[string]string
	Symbols   map[string]lsp.Symbol
}

func NewState() State {
	return State{
		Documents: map[string]string{},
		Symbols:   map[string]lsp.Symbol{},
	}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

// TODO: Where we will do Type
func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	symbol, err := GetSymbolAtPosition(uri, position)
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

	contents := fmt.Sprintf("Symbol: %s\nType: %s\nDataType: %s", symbol.Name, symbol.Type, symbol.DataType)

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

func GetSymbolAtPosition(uri string, position lsp.Position) (lsp.Symbol, error) {
	return lsp.Symbol{
		Name:     "example-symbol",
		Type:     "variable",
		DataType: "int",
	}, nil
}
