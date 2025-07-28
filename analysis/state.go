package analysis

import "vlogcc-lsp/lsp"

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

//TODO: Where we will do Type
func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: "Hello VLOGCC language",
		},
	}
}
