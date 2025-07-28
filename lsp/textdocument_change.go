package lsp

type TextDocumentDidChangeNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}
type DidChangeTextDocumentParams struct {
	TextDocument VersionTextDocumentIdentifier `json:"textDocument"`
	ContentChanges []TextDocumentChangeEvent `json:"contentChanges"`
}

type TextDocumentChangeEvent struct {
	Text string `json:"text"`
}
