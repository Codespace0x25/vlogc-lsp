package lsp

type TextDocumentItem struct {
	URI string `json:"uri"`
	LanguageId string `json:"languageid"`
	Version int `json:"version"`
	Text string `json:"text"`
}
